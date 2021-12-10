package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		email := jsoniter.Get(scanner.Bytes(), "Email").ToString()

		if strings.HasSuffix(email, "."+domain) {
			parts := strings.Split(email, "@")
			key := strings.ToLower(parts[1])
			result[key]++
		}
	}

	return result, nil
}
