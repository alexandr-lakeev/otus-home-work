package logger

import (
	"fmt"
	"testing"
)

func TestLogger(t *testing.T) {
	logg, err := New("INFO")

	fmt.Println(logg, err)
}
