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
		sb     strings.Builder
		last   string
		escape bool
	)

	for _, r := range s {
		if unicode.IsDigit(r) && !escape {
			n, err := strconv.Atoi(string(r))
			if err != nil {
				return "", err
			}

			if last == "" {
				return "", ErrInvalidString
			}

			sb.WriteString(strings.Repeat(last, n))
			last = ""
			continue
		}

		if r == '\\' && !escape {
			escape = true
			continue
		}

		sb.WriteString(last)
		last = string(r)
		escape = false
	}

	sb.WriteString(last)

	return sb.String(), nil
}
