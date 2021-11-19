package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	c := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec

	c.Env = createCmdEnv(env)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout

	if err := c.Run(); err != nil {
		return 1
	}

	return 0
}

func createCmdEnv(env Environment) []string {
	cmdEnv := os.Environ()

	for name, value := range env {
		cmdEnv = append(cmdEnv, name+"="+value.Value)
	}

	return cmdEnv
}
