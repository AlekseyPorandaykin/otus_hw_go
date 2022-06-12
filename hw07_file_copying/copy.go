package main

import (
	"errors"
	"io"
	"math"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrWriteToFile           = errors.New("error write to file")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	readFile, sizeFile, errFile := getReadFile(fromPath, offset)
	if errFile != nil {
		return ErrUnsupportedFile
	}

	writeFile, errWrite := getWriteFile(toPath)
	if errWrite != nil {
		return ErrUnsupportedFile
	}
	defer writeFile.Close()
	if err := copyFile(readFile, writeFile, limit, sizeFile); err != nil {
		return err
	}
	return nil
}

func getReadFile(fromPath string, offset int64) (*os.File, int64, error) {
	file, errFile := os.Open(fromPath)
	if errFile != nil {
		return nil, 0, errFile
	}
	stat, errStat := file.Stat()
	if errStat != nil {
		return nil, 0, errStat
	}
	if offset >= stat.Size() || stat.Size() == 0 {
		return nil, 0, ErrOffsetExceedsFileSize
	}
	if _, errSeek := file.Seek(offset, io.SeekStart); errSeek != nil {
		return nil, 0, errSeek
	}
	sizeFile := stat.Size() - offset

	return file, sizeFile, nil
}

func getWriteFile(toPath string) (*os.File, error) {
	file, errCreate := os.Create(toPath)
	if errCreate != nil {
		return nil, errCreate
	}

	return file, nil
}

func copyFile(reader io.Reader, writer io.Writer, limit, sizeFile int64) error {
	if limit > 0 {
		reader = io.LimitReader(reader, limit)
		sizeFile = limit
	}
	sizeBuf := 1 << 8
	bar := pb.StartNew(int(math.Ceil(float64(sizeFile) / float64(sizeBuf))))
	buf := make([]byte, sizeBuf)
	for {
		rn, errRead := reader.Read(buf)
		if io.EOF == errRead {
			break
		}
		if errRead != nil {
			return ErrWriteToFile
		}
		wn, errWrite := writer.Write(buf[0:rn])
		if errWrite != nil || wn < rn {
			return ErrWriteToFile
		}
		bar.Increment()
	}
	bar.Finish()

	return nil
}
