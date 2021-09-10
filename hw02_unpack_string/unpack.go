package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var (
		sb      strings.Builder
		last    string
		escaped bool
	)

	for _, r := range s {
		switch {
		case !escaped && r == '\\':
			escaped = true
		case !escaped && unicode.IsDigit(r):
			n, err := strconv.Atoi(string(r))
			if err != nil {
				return "", err
			}
			if last == "" {
				return "", ErrInvalidString
			}
			sb.WriteString(strings.Repeat(last, n))
			last = ""
			escaped = false
		default:
			sb.WriteString(last)
			last = string(r)
			escaped = false
		}
	}

	sb.WriteString(last)

	return sb.String(), nil
}
