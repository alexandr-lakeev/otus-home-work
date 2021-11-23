package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("Read env", func(t *testing.T) {
		env, err := ReadDir("testdata/env")

		require.NoError(t, err)

		require.False(t, env["BAR"].NeedRemove)
		require.False(t, env["EMPTY"].NeedRemove)
		require.False(t, env["FOO"].NeedRemove)
		require.False(t, env["HELLO"].NeedRemove)
		require.True(t, env["UNSET"].NeedRemove)

		require.Equal(t, "bar", env["BAR"].Value)
		require.Equal(t, "", env["EMPTY"].Value)
		require.Equal(t, "   foo\nwith new line", env["FOO"].Value)
		require.Equal(t, "\"hello\"", env["HELLO"].Value)
		require.Equal(t, "", env["UNSET"].Value)
	})
}
