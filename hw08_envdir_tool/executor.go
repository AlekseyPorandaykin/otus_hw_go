package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	returnCode = 0
	setEnvs(env)
	if len(cmd) < 1 {
		return 1
	}
	cmdName := cmd[0]
	args := cmd[1:]
	if len(cmdName) < 1 {
		return 1
	}
	wcExel := exec.Command(cmdName, args...)
	wcExel.Stdout = os.Stdout
	wcExel.Stderr = os.Stderr
	wcExel.Stdin = os.Stdin
	if err := wcExel.Run(); err != nil {
		returnCode = 1
	}

	return returnCode
}

func setEnvs(env Environment) {
	for name, envVal := range env {
		os.Unsetenv(name)
		if !envVal.NeedRemove {
			os.Setenv(name, envVal.Value)
		}
	}
}
