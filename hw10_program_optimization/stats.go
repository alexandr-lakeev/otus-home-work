package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)

	domain = "." + domain
	scanner := bufio.NewScanner(r)

	var sb strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		collect := false
		for _, r := range line {
			if r == '@' {
				collect = true
				continue
			}

			if collect {
				if r == '"' {
					break
				} else {
					sb.WriteRune(r)
				}
			}
		}

		if strings.HasSuffix(sb.String(), domain) {
			key := strings.ToLower(sb.String())
			result[key]++
		}

		sb.Reset()
	}

	return result, nil
}
