package main

import (
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	c := exec.Command(cmd[0], cmd[1:]...)

	c.Env = createCmdEnv(env)
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin

	err := c.Run()
	if err != nil {
		log.Println(err)
		return 1
	}

	return 0
}

func createCmdEnv(env Environment) []string {
	var cmdEnv = os.Environ()

	for name, value := range env {
		cmdEnv = append(cmdEnv, name+"="+value.Value)
	}

	return cmdEnv
}
