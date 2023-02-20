Feature: Send event

  Scenario: Create new event 1
    When I send "POST" request to "/events" with data
      """
		{
            "title": "Test  integration title 1",
            "dateTimeStart": "2022-02-12 17:20:06",
            "dateTimeEnd":"2022-02-12 17:24:05",
            "description":"Test description",
            "createdBy":1,
            "remindFrom": "2022-02-10 11:24:05"
        }
		"""
    Then I response code should be 201
    And Response data has code "Event created success"

  Scenario: Error create exist event
    When I send "POST" request to "/events" with data
      """
		{
            "title": "Test  integration title 1",
            "dateTimeStart": "2022-02-12 17:20:06",
            "dateTimeEnd":"2022-02-12 17:24:05",
            "description":"Test description",
            "createdBy":1,
            "remindFrom": "2022-02-10 11:24:05"
        }
		"""
    Then I response code should be 500
    And Response data has code "Error create event"

  Scenario: Create new event 2
    When I send "POST" request to "/events" with data
      """
		{
            "title": "Test  integration title 2",
            "dateTimeStart": "2022-02-16 17:20:06",
            "dateTimeEnd":"2022-02-16 14:24:05",
            "description":"Test description",
            "createdBy":1,
            "remindFrom": "2022-02-16 11:24:05"
        }
		"""
    Then I response code should be 201
    And Response data has code "Event created success"

  Scenario: Create new event 3
    When I send "POST" request to "/events" with data
      """
		{
            "title": "Test  integration title 3",
            "dateTimeStart": "2022-02-21 17:20:06",
            "dateTimeEnd":"2022-02-21 17:24:05",
            "description":"Test description",
            "createdBy":1,
            "remindFrom": "2022-02-21 11:24:05"
        }
		"""
    Then I response code should be 201
    And Response data has code "Event created success"

  Scenario: Create new event 4 without remind
    When I send "POST" request to "/events" with data
      """
		{
            "title": "Test  integration title 4",
            "dateTimeStart": "2022-02-12 18:30:06",
            "dateTimeEnd":"2022-02-12 18:24:05",
            "description":"Test description",
            "createdBy":1,
            "remindFrom": ""
        }
		"""
    Then I response code should be 201
    And Response data has code "Event created success"

  Scenario: Get events on day
    When I send "GET" request to "/events/day/2022-02-12"
    Then I response code should be 200
    And Response data has code "Events read success"
    And Has event with title "Test  integration title 1"
    And Has event with title "Test  integration title 4"
    And Has not event with title "Test  integration title 2"
    And Has not event with title "Test  integration title 3"


  Scenario: Get events on week
    When I send "GET" request to "/events/week/2022-02-12"
    Then I response code should be 200
    And Response data has code "Events read success"
    And Has event with title "Test  integration title 1"
    And Has event with title "Test  integration title 2"
    And Has event with title "Test  integration title 4"
    And Has not event with title "Test  integration title 3"

  Scenario: Get events on month
    When I send "GET" request to "/events/month/2022-02-12"
    Then I response code should be 200
    And Response data has code "Events read success"
    And Has event with title "Test  integration title 1"
    And Has event with title "Test  integration title 2"
    And Has event with title "Test  integration title 3"
    And Has event with title "Test  integration title 4"

  Scenario: Check send notification
    When Wait "12s" when scheduler send all notification
    Then Find in log event with title "Test  integration title 1"
    And Find in log event with title "Test  integration title 2"
    And Find in log event with title "Test  integration title 3"
    And Not find in log event with title "Test  integration title 4"

