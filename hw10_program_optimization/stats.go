package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)

	domain = "." + domain
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		userDomain := getDomainFromLine(line)

		if strings.HasSuffix(userDomain, domain) {
			key := strings.ToLower(userDomain)
			result[key]++
		}
	}

	return result, nil
}

func getDomainFromLine(line string) string {
	var sb strings.Builder

	startIndex := strings.Index(line, "\"Email\":\"")
	if startIndex == -1 {
		return ""
	}

	collect := false
	for n, r := range line {
		if n < startIndex {
			continue
		}

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

	return sb.String()
}
