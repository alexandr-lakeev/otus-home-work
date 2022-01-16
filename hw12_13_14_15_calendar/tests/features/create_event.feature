Feature: Event creating
	Scenario: Event create successfully
		When I send "POST" request to "http://backend:8080/api/calendar/v1/events" with "application/json" data:
		"""
		{
			"first_name": "otus",
			"email": "otus@otus.ru",
			"age": 27
		}
		"""
		Then The response code should be 200