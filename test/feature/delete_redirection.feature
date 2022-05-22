# language: en

Feature: Delete a short link

    Scenario: The short link exists
        Given that I got a short link for http://www.google.com
        When I command the system to delete the short link
        Then the system confirms that the operation was successfully executed
    
    Scenario: The short link does not exist
        Given that I got a short link that does not exist
        When I command the system to delete the short link
        Then the system signals that the short link does not exist