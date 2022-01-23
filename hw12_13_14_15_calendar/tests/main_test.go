package tests

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

const delay = 10 * time.Second

func TestMain(m *testing.M) {
	log.Printf("wait %s for service availability...", delay)
	time.Sleep(delay)

	suite := godog.TestSuite{
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			FeatureContext(s)
		},
		Options: &godog.Options{
			Format:    "pretty",
			Paths:     []string{"features"},
			Randomize: 0,
		},
	}

	status := suite.Run()

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
