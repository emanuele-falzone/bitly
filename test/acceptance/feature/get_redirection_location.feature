# language: en

Feature: Navigate to a short URL

    Scenario: Navigate to a previouly created short URL
        Given that I got a short link for http://www.google.com
        When I navigate to the short link
        Then the system redirects me to http://www.google.com
    
    Scenario: Navigate to a short URL taht does not exists
        Given that I got a short link that does not exist
        When I navigate to the short link
        Then the system signals that the short link does not exist