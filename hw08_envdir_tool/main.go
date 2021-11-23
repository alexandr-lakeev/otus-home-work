package main

import (
	"errors"
	"log"
	"os"
)

var ErrEnoughArguments = errors.New("not enough arguments")

func main() {
	args := os.Args

	if len(args) < 3 {
		log.Fatal(ErrEnoughArguments)
	}

	dir, args := args[1], args[2:]

	env, err := ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(RunCmd(args, env))
}
