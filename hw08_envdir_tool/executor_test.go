package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	tests := []struct {
		nameTest         string
		nameEnv          string
		needRemove       bool
		envValue         string
		expectedValueEnv string
	}{
		{
			nameTest:         "Test set env value",
			nameEnv:          "TEST_CORRECT_ENV",
			needRemove:       false,
			envValue:         "some data",
			expectedValueEnv: "some data",
		},
		{
			nameTest:         "Test replace empty env value",
			nameEnv:          "TEST_ENV_ONE",
			needRemove:       false,
			envValue:         "",
			expectedValueEnv: "",
		},
		{
			nameTest:         "Test remove env",
			nameEnv:          "TEST_ENV_ONE",
			needRemove:       true,
			envValue:         "",
			expectedValueEnv: "",
		},
	}
	os.Setenv("TEST_ENV_ONE", "value one")
	for _, testData := range tests {
		t.Run(testData.nameTest, func(t *testing.T) {
			envs := map[string]EnvValue{}
			envs[testData.nameEnv] = EnvValue{Value: testData.envValue, NeedRemove: testData.needRemove}
			RunCmd([]string{"cd", "testdata"}, envs)
			require.Equal(t, testData.expectedValueEnv, os.Getenv(testData.nameEnv))
		})
	}
}
