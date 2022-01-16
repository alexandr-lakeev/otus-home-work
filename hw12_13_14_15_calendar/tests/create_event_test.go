package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v16"
	"github.com/google/uuid"
)

type createEventTest struct {
	userID             uuid.UUID
	responseStatusCode int
	responseBody       []byte
}

func (test *createEventTest) createUserId(*messages.Pickle) {
	test.userID = uuid.New()
}

func (test *createEventTest) iSendRequestToWithData(httpMethod, addr, contentType string, data *godog.DocString) error {
	var (
		r   *http.Response
		err error
	)

	switch httpMethod {
	case http.MethodPost:
		replacer := strings.NewReplacer("\n", "", "\t", "")
		cleanJson := replacer.Replace(data.Content)

		url, err := url.Parse(addr)
		if err != nil {
			return err
		}
		body := ioutil.NopCloser(strings.NewReader(cleanJson))

		var client = &http.Client{}
		request := &http.Request{
			Method: httpMethod,
			URL:    url,
			Body:   body,
			Header: map[string][]string{
				"Content-Type": {contentType},
				"X-User-Id":    {test.userID.String()},
			},
		}
		r, err = client.Do(request)
	default:
		err := fmt.Errorf("unknown method: %s", httpMethod)
		if err != nil {
			return err
		}
	}

	test.responseStatusCode = r.StatusCode
	test.responseBody, err = ioutil.ReadAll(r.Body)

	return err
}

func (test *createEventTest) theResponseCodeShouldBe(code int) error {
	if test.responseStatusCode != code {
		return fmt.Errorf("unexpected status code: %d != %d", test.responseStatusCode, code)
	}
	return nil
}

func (test *createEventTest) theErrorShouldBeInTheField(errText, errField string) error {
	response := struct {
		Data   map[string]interface{} `json:"data"`
		Errors map[string]string      `json:"errors"`
	}{}

	err := json.Unmarshal([]byte(test.responseBody), &response)
	if err != nil {
		return err
	}

	if response.Errors[errField] != errText {
		return fmt.Errorf("unexpected json: %s != %s", response.Errors[errField], errText)
	}
	return nil
}

func FeatureContext(s *godog.ScenarioContext) {
	test := new(createEventTest)

	s.BeforeScenario(test.createUserId)

	s.Step(`^I send "([^"]*)" request to "([^"]*)" with "([^"]*)" data:$`, test.iSendRequestToWithData)
	s.Step(`^The response code should be (\d+)$`, test.theResponseCodeShouldBe)
	s.Step(`^The error should be "([^"]*)" in The field "([^"]*)"$`, test.theErrorShouldBeInTheField)
}
