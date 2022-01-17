Feature: Get Events list
	Scenario: Events list for a day
		When I send "GET" request to "http://backend:8080/api/calendar/v1/events?from=2021-12-31T00:00:00Z&to=2021-12-31T23:59:59Z"
		Then The response code should be 200
		And The json should contain 1 event

	Scenario: Events list for a week
		When I send "GET" request to "http://backend:8080/api/calendar/v1/events?from=2021-12-27T00:00:00Z&to=2022-01-02T23:59:59Z"
		Then The response code should be 200
		And The json should contain 3 events

	Scenario: Events list for a month
		When I send "GET" request to "http://backend:8080/api/calendar/v1/events?from=2021-12-01T00:00:00Z&to=2021-12-31T23:59:59Z"
		Then The response code should be 200
		And The json should contain 5 events
