Feature: Send notificaion
	Scenario: Notification event is received
	  When I send "POST" request to "http://backend:8080/api/calendar/v1/events" with "application/json" data:
		"""
    {
      "title":"Otus hw15_calendar have done",
      "description":"Event description",
      "date":"2021-12-31T23:59:59Z",
      "duration":30
    }
		"""
		Then The response code should be 200
		And I receive event with title "Otus hw15_calendar have done"
