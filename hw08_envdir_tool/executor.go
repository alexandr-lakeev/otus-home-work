package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	c := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec

	cleanEnv(env)

	c.Env = createCmdEnv(env)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		var exitError *exec.ExitError

		if errors.As(err, &exitError) {
			returnCode = exitError.ExitCode()
		}
	}

	return
}

func cleanEnv(env Environment) {
	for name := range env {
		os.Unsetenv(name)
	}
}

func createCmdEnv(env Environment) []string {
	cmdEnv := os.Environ()

	for name, value := range env {
		if !value.NeedRemove {
			cmdEnv = append(cmdEnv, name+"="+value.Value)
		}
	}

	return cmdEnv
}
