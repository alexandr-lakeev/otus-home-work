Feature: Event creating
	Scenario: Event create successfully
		When I send "POST" request to "http://backend:8080/api/calendar/v1/events" with "application/json" data:
		"""
    {
      "title":"Dummy event",
      "description":"Event description",
      "date":"2021-12-31T23:59:59Z",
      "duration":30
    }
		"""
		Then The response code should be 200

  Scenario: Get error when trying to create an event on the same date
  	When I send "POST" request to "http://backend:8080/api/calendar/v1/events" with "application/json" data:
		"""
    {
      "title":"Dummy event",
      "description":"Event description",
      "date":"2021-12-31T23:59:59Z",
      "duration":30
    }
		"""
		Then The response code should be 200
		When I send "POST" request to "http://backend:8080/api/calendar/v1/events" with "application/json" data:
		"""
    {
      "title":"Dummy event 2",
      "description":"Event description 2",
      "date":"2021-12-31T23:59:59Z",
      "duration":30
    }
		"""
		Then The response code should be 400
    And The error should be "date is busy by other event" in The field "error"

  Scenario: Get error when trying to create an event with wrong data
  	When I send "POST" request to "http://backend:8080/api/calendar/v1/events" with "application/json" data:
		"""
    {}
		"""
		Then The response code should be 400
    And The error should be "this field is required" in The field "Title"
    And The error should be "this field is required" in The field "Date"
    And The error should be "this field is required" in The field "Duration"
