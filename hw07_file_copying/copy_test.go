package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

//nolint
var testData = `Go
Documents
Packages
The Project
Help
Blog
Play
Search

Getting Started
Install the Go tools
Test your installation
Installing extra Go versions
Uninstalling Go
Getting help
Download the Go distribution
Download Go
Click here to visit the downloads page
Official binary distributions are available for the FreeBSD (release 10-STABLE and above), Linux, macOS (10.10 and above), and Windows operating systems and the 32-bit (386) and 64-bit (amd64) x86 processor architectures.

If a binary distribution is not available for your combination of operating system and architecture, try installing from source or installing gccgo instead of gc.`

func TestCopy(t *testing.T) {
	tmp := os.TempDir()

	input := tmp + "/input.txt"
	out := tmp + "/out.txt"

	err := os.WriteFile(input, []byte(testData), 0775)
	require.NoError(t, err)
	defer os.Remove(input)

	t.Run("Copy", func(t *testing.T) {
		copyParams := []struct {
			name     string
			offset   int64
			limit    int64
			expected string
		}{
			{"full", 0, 0, testData},
			{"limit exceeds file size", 0, 999999, testData},
			{"offset 0 limit 12", 0, 12, "Go\nDocuments"},
			{"offset 3 limit 9", 3, 9, "Documents"},
			{"offset 626 to the end", 626, 0, "instead of gc."},
		}

		for _, cp := range copyParams {
			cp := cp
			t.Run(cp.name, func(t *testing.T) {
				err := Copy(input, out, cp.offset, cp.limit)

				require.NoError(t, err)

				defer os.Remove(out)

				resultFile, _ := os.Open(out)
				resultBytes, _ := ioutil.ReadAll(resultFile)

				require.Equal(t, cp.expected, string(resultBytes))
			})
		}
	})

	t.Run("Offset exceeds file size", func(t *testing.T) {
		err := Copy(input, out, 999999, 0)

		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("Unsupported file", func(t *testing.T) {
		err := Copy("/dev/urandom", out, 0, 0)

		require.ErrorIs(t, err, ErrUnsupportedFile)
	})
}
