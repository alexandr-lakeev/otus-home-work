package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v16"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

const amqpURL = "amqp://user:password@rabbit:5672/"

const (
	queueName    = "test.integration"
	exchangeName = "test"
)

type calendarTest struct {
	conn          *amqp.Connection
	ch            *amqp.Channel
	messages      [][]byte
	messagesMutex *sync.RWMutex
	stopSignal    chan struct{}

	userID             uuid.UUID
	responseStatusCode int
	responseBody       []byte
}

func (test *calendarTest) startConsuming(*messages.Pickle) {
	test.messages = make([][]byte, 0)
	test.messagesMutex = new(sync.RWMutex)
	test.stopSignal = make(chan struct{})

	var err error

	test.conn, err = amqp.Dial(amqpURL)
	if err != nil {
		panic(err)
	}

	test.ch, err = test.conn.Channel()
	if err != nil {
		panic(err)
	}

	_, err = test.ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	err = test.ch.QueueBind(queueName, "", exchangeName, false, nil)
	if err != nil {
		panic(err)
	}

	events, err := test.ch.Consume(queueName, "", true, true, false, false, nil)
	if err != nil {
		panic(err)
	}

	go func(stop <-chan struct{}) {
		for {
			select {
			case <-stop:
				return
			case event := <-events:
				test.messagesMutex.Lock()
				test.messages = append(test.messages, event.Body)
				test.messagesMutex.Unlock()
			}
		}
	}(test.stopSignal)
}

func (test *calendarTest) stopConsuming(*messages.Pickle, error) {
	test.stopSignal <- struct{}{}

	if err := test.ch.Close(); err != nil {
		panic(err)
	}

	if err := test.conn.Close(); err != nil {
		panic(err)
	}

	test.messages = nil
}

func (test *calendarTest) createUserId(*messages.Pickle) {
	test.userID = uuid.New()
}

func (test *calendarTest) addTestEvents(*messages.Pickle) {
	testEventsJson := []string{
		`{
      "title":"Event #1",
      "date":"2021-12-02T13:00:00Z",
      "duration":60
    }`,
		`{
      "title":"Event #2",
      "date":"2021-12-15T14:00:00Z",
      "duration":60
    }`,
		`{
      "title":"Event #3",
      "date":"2021-12-28T15:00:00Z",
      "duration":60
    }`,
		`{
      "title":"Event #4",
      "date":"2021-12-30T16:00:00Z",
      "duration":60
    }`,
		`{
      "title":"Event #5",
      "date":"2021-12-31T18:00:00Z",
      "duration":60
    }`,
	}

	url, err := url.Parse("http://backend:8080/api/calendar/v1/events")
	if err != nil {
		panic(err)
	}

	for _, eventJson := range testEventsJson {
		body := ioutil.NopCloser(strings.NewReader(eventJson))

		client := http.Client{
			Timeout: 1 * time.Second,
		}
		request := &http.Request{
			Method: http.MethodPost,
			URL:    url,
			Body:   body,
			Header: map[string][]string{
				"Content-Type": {"application/json"},
				"X-User-Id":    {test.userID.String()},
			},
		}
		r, err := client.Do(request)
		if err != nil {
			panic(err)
		}

		if r.StatusCode != http.StatusOK {
			panic(fmt.Errorf("test event not created, status code %d", r.StatusCode))
		}
	}
}

func (test *calendarTest) iSendRequestToWithData(httpMethod, addr, contentType string, data *godog.DocString) error {
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

		client := http.Client{}
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

func (test *calendarTest) iSendRequestTo(httpMethod, addr string) error {
	var (
		r   *http.Response
		err error
	)

	switch httpMethod {
	case http.MethodGet:
		url, err := url.Parse(addr)
		if err != nil {
			return err
		}

		client := http.Client{
			Timeout: 1 * time.Second,
		}
		request := &http.Request{
			Method: httpMethod,
			URL:    url,
			Header: map[string][]string{
				"X-User-Id": {test.userID.String()},
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

func (test *calendarTest) theResponseCodeShouldBe(code int) error {
	if test.responseStatusCode != code {
		return fmt.Errorf("unexpected status code: %d != %d", test.responseStatusCode, code)
	}
	return nil
}

func (test *calendarTest) theErrorShouldBeInTheField(errText, errField string) error {
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

func (test *calendarTest) theJsonShouldContainEvents(eventsCount int) error {
	response := struct {
		Data []interface{} `json:"data"`
	}{}

	err := json.Unmarshal([]byte(test.responseBody), &response)
	if err != nil {
		return err
	}

	if len(response.Data) != eventsCount {
		return fmt.Errorf("unexpected events count: %d != %d", len(response.Data), eventsCount)
	}
	return nil
}

func (test *calendarTest) iReceiveEventWithTitle(eventTitle string) error {
	time.Sleep(10 * time.Second)

	test.messagesMutex.RLock()
	defer test.messagesMutex.RUnlock()

	var event map[string]string

	for _, msg := range test.messages {
		event = make(map[string]string)
		json.Unmarshal(msg, &event)

		if event["Title"] == eventTitle {
			return nil
		}
	}

	return fmt.Errorf("event with title '%s' was not found in %s", eventTitle, test.messages)
}

func FeatureContext(s *godog.ScenarioContext) {
	test := new(calendarTest)

	s.BeforeScenario(test.startConsuming)
	s.BeforeScenario(test.createUserId)
	s.BeforeScenario(test.addTestEvents)

	s.Step(`^I send "([^"]*)" request to "([^"]*)" with "([^"]*)" data:$`, test.iSendRequestToWithData)
	s.Step(`^I send "([^"]*)" request to "([^"]*)"$`, test.iSendRequestTo)
	s.Step(`^The response code should be (\d+)$`, test.theResponseCodeShouldBe)
	s.Step(`^The error should be "([^"]*)" in The field "([^"]*)"$`, test.theErrorShouldBeInTheField)
	s.Step(`^The json should contain (\d+) event[s]?$`, test.theJsonShouldContainEvents)
	s.Step(`^I receive event with title "([^"]*)"$`, test.iReceiveEventWithTitle)

	s.AfterScenario(test.stopConsuming)
}
