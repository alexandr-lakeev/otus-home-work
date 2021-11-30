package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("Copy file", func(t *testing.T) {
		returnCode := RunCmd([]string{"cp", "./testdata/env/BAR", "/tmp"}, nil)

		require.Zero(t, returnCode)

		data, err := os.ReadFile("/tmp/BAR")

		require.NoError(t, err)

		os.Remove("/tmp/BAR")

		require.Equal(t, "bar\nPLEASE IGNORE SECOND LINE\n", string(data))
	})

	t.Run("Non zero exit code", func(t *testing.T) {
		returnCode := RunCmd([]string{"cp", "./testdata/env/NOFILE", "/tmp/NOFILE"}, nil)

		require.Equal(t, 1, returnCode)
	})

	t.Run("Environment", func(t *testing.T) {
		fakeStdoutPath := "/tmp/stdout"
		fakeStdout, _ := os.Create(fakeStdoutPath)

		os.Stdout = fakeStdout

		env := Environment{
			"FOO": EnvValue{
				Value:      "foo",
				NeedRemove: false,
			},
			"BAR": EnvValue{
				Value:      "bar",
				NeedRemove: false,
			},
		}

		RunCmd([]string{"printenv", "FOO", "BAR"}, env)

		data, err := os.ReadFile(fakeStdoutPath)

		require.NoError(t, err)

		os.Remove(fakeStdoutPath)

		require.Equal(t, "foo\nbar\n", string(data))
	})
}
