package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"strings"
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
	filesInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	envs := make(Environment)

	for _, fi := range filesInfo {
		name := fi.Name()

		envs[name], err = readEnvValue(fi, dir)

		if err != nil {
			return nil, err
		}
	}

	return envs, nil
}

func readEnvValue(fileInfo os.FileInfo, dir string) (EnvValue, error) {
	var value EnvValue

	if fileInfo.Size() == 0 {
		value.NeedRemove = true
		return value, nil
	}

	fullpath := path.Join(dir, fileInfo.Name())

	file, err := os.Open(fullpath)
	if err != nil {
		return value, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()

	value.Value = strings.TrimRight(scanner.Text(), " \n\t")
	value.Value = string(bytes.ReplaceAll([]byte(value.Value), []byte{0x00}, []byte("\n")))

	return value, nil
}
