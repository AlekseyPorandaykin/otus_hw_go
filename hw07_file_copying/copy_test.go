package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type testStruct struct {
	testName string
	errorStr string
	from     string
	to       string
	offset   int64
	limit    int64
}

func TestCopyErrors(t *testing.T) {
	testsWithoutCreateFile := []testStruct{
		{
			testName: "Test failed file",
			errorStr: "unsupported file",
			from:     "testdata/dummy_from.txt",
			to:       "dummy_to.txt",
		},
		{
			testName: "Test failed read unlimited file",
			errorStr: "unsupported file",
			from:     "/dev/urandom",
			to:       "error_out.txt",
		},
		{
			testName: "Test write to incorrect file",
			errorStr: "unsupported file",
			from:     "testdata/out_offset6000_limit1000.txt",
			to:       "",
		},
		{
			testName: "Test copy empty file",
			errorStr: "unsupported file",
			from:     "testdata/empty.txt",
		},
	}

	for _, testData := range testsWithoutCreateFile {
		td := testData
		t.Run(td.testName, func(t *testing.T) {
			errCopy := Copy(td.from, td.to, td.offset, td.limit)
			require.Equal(t, td.errorStr, errCopy.Error())
			require.NoFileExists(t, td.to)
		})
	}
	testsWithCreateFile := []testStruct{
		{
			testName: "Test offset more length data",
			errorStr: "unsupported file",
			from:     "testdata/out_offset6000_limit1000.txt",
			offset:   1001,
		},
	}
	for _, testData := range testsWithCreateFile {
		td := testData
		t.Run(td.testName, func(t *testing.T) {
			file, errCreateFile := os.CreateTemp("", "copy*.txt")
			file.Name()
			if errCreateFile != nil {
				t.Fatal("Error create tmp file")
			}
			defer func() {
				errDeleteFile := os.Remove(file.Name())
				if errDeleteFile != nil {
					t.Fatalf("Error delete temp file: %s", file.Name())
				}
			}()
			errCopy := Copy(td.from, file.Name(), td.offset, td.limit)
			require.Equal(t, td.errorStr, errCopy.Error())
		})
	}
}

func TestCopy(t *testing.T) {
	t.Run("Copy full file without limit and offset", func(t *testing.T) {
		file, errCreateFile := os.CreateTemp("", "copy*.txt")
		file.Name()
		if errCreateFile != nil {
			t.Fatal("Error create tmp file")
		}
		defer func() {
			errDeleteFile := os.Remove(file.Name())
			if errDeleteFile != nil {
				t.Fatalf("Error delete temp file: %s", file.Name())
			}
		}()
		_ = Copy("testData/input.txt", file.Name(), 0, 0)
		require.Equal(t, getSizeFile("testData/input.txt"), getSizeFile(file.Name()))
	})
	t.Run("Copy full file with big limit", func(t *testing.T) {
		file, errCreateFile := os.CreateTemp("", "copy*.txt")
		file.Name()
		if errCreateFile != nil {
			t.Fatal("Error create tmp file")
		}
		defer func() {
			errDeleteFile := os.Remove(file.Name())
			if errDeleteFile != nil {
				t.Fatalf("Error delete temp file: %s", file.Name())
			}
		}()
		_ = Copy("testData/input.txt", file.Name(), 0, 10000000)
		require.Equal(t, getSizeFile("testData/input.txt"), getSizeFile(file.Name()))
	})
}

func getSizeFile(pathToFile string) int64 {
	file, _ := os.Open(pathToFile)
	info, err := file.Stat()
	if err != nil {
		return 0
	}
	return info.Size()
}
