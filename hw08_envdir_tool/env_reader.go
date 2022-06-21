package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
)

var (
	ErrReadEnvFromDir      = errors.New("error read env from directory")
	ErrReadDataFromEnvFile = errors.New("error read data from env file")

	EOF          = byte(10)
	SpaceStr     = byte(32)
	TabStr       = byte(9)
	TerminalNull = byte(0)
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	envs := map[string]EnvValue{}
	fileInfos, errReadEnvs := ioutil.ReadDir(dir)
	if errReadEnvs != nil {
		return nil, ErrReadEnvFromDir
	}
	for _, fileInfo := range fileInfos {
		if !isCorrectFileName(fileInfo.Name()) {
			continue
		}
		dataFromFile, errReadData := os.ReadFile(fmt.Sprintf("%s/%s", dir, fileInfo.Name()))
		if errReadData != nil {
			return nil, ErrReadDataFromEnvFile
		}
		envs[fileInfo.Name()] = generateEnvValue(dataFromFile)
	}
	return envs, nil
}

func generateEnvValue(dataFromFile []byte) EnvValue {
	var bufStr strings.Builder
	for _, symbol := range dataFromFile {
		if symbol == EOF {
			break
		}
		bufStr.WriteByte(symbol)
	}
	return EnvValue{
		Value:      trimLastEmptySymbols(replaceTerminalNull(bufStr.String())),
		NeedRemove: len(dataFromFile) < 1,
	}
}

func isCorrectFileName(fileName string) bool {
	return !strings.Contains(fileName, "=") && isUpperStr(fileName)
}

func isUpperStr(str string) bool {
	for _, sym := range str {
		if !unicode.IsLetter(sym) {
			continue
		}
		if !unicode.IsUpper(sym) {
			return false
		}
	}
	return true
}

func replaceTerminalNull(str string) string {
	return strings.ReplaceAll(str, string(TerminalNull), string(EOF))
}

func trimLastEmptySymbols(str string) string {
	byteStr := []byte(str)
	for {
		if len(byteStr) < 1 {
			break
		}
		lastByte := byteStr[len(byteStr)-1]
		if lastByte != SpaceStr && lastByte != TabStr {
			break
		}
		byteStr = byteStr[:len(byteStr)-1]
	}

	return string(byteStr)
}
