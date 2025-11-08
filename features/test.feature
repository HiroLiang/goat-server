Feature: Standard feature test

  Scenario: Show how to use api test context
    When I send a "GET" request to "/api/test"
    Then the response should contain {"result": "Test OK"}