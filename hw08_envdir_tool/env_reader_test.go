package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	err := os.Mkdir("./testdata/test", 0o777)
	if err != nil {
		t.Fatal("Error create dir")
	}
	defer os.Remove("./testdata/test")
	tests := []struct {
		nameTest      string
		nameEnv       string
		expectedValue string
		dataFile      string
		needRemove    bool
	}{
		{
			nameTest:      "Test empty symbols",
			nameEnv:       "TEST_EMPTY_SYMBOLS",
			expectedValue: "data",
			dataFile:      "data\t \t \t \nextra data",
			needRemove:    false,
		},
		{
			nameTest:      "Test incorrect name file, incorrect symbol",
			nameEnv:       "TEST_INCORRECT_=FILE",
			expectedValue: "",
			dataFile:      "empty data",
			needRemove:    false,
		},
		{
			nameTest:      "Test replace terminal null",
			nameEnv:       "TEST_REPLACE_TERMINAL_NULL",
			expectedValue: "first part\nsecond part\nlast part",
			dataFile:      "first part\u0000second part\u0000last part\t\nextra line",
			needRemove:    false,
		},
		{
			nameTest:      "Test delete env value",
			nameEnv:       "TEST_DELETE_ENV",
			expectedValue: "",
			dataFile:      "",
			needRemove:    true,
		},
		{
			nameTest:      "Test set empty env value",
			nameEnv:       "TEST_SET_EMPTY_VALUE",
			expectedValue: "",
			dataFile:      "\t \ntest",
			needRemove:    false,
		},
	}

	for _, testData := range tests {
		testData := testData
		dir := "./testdata/test"
		t.Run(testData.nameTest, func(t *testing.T) {
			file, errFile := os.Create(fmt.Sprintf("%s/%s", dir, testData.nameEnv))
			if errFile != nil {
				t.Fatal("Error create tmp file")
			}
			file.Write([]byte(testData.dataFile))
			file.Close()
			envs, _ := ReadDir("./testdata/test")
			envValue := envs[testData.nameEnv]
			require.Equal(t, testData.expectedValue, envValue.Value)
			require.Equal(t, testData.needRemove, envValue.NeedRemove)
			defer func() {
				os.Remove(file.Name())
			}()
		})
	}
}

func TestErrorReadDir(t *testing.T) {
	t.Run("Incorrect dir", func(t *testing.T) {
		envs, err := ReadDir("./dummyDir/incorrect_dir")
		require.Nil(t, envs)
		require.Equal(t, ErrReadEnvFromDir, err)
	})
}
