package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	json "github.com/json-iterator/go"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)

	domain = "." + domain
	scanner := bufio.NewScanner(r)

	user := &User{}
	for scanner.Scan() {
		*user = User{}
		if err := json.Unmarshal([]byte(scanner.Text()), user); err != nil {
			return nil, err
		}

		if strings.HasSuffix(user.Email, domain) {
			key := strings.ToLower(strings.Split(user.Email, "@")[1])
			result[key]++
		}
	}

	return result, nil
}
