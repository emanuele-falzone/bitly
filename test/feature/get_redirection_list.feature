# language: en

Feature: Get redirection list

    Scenario: Navigate to a short URL taht does not exists
        Given that I got a short link for http://www.google.com
        And that I got a short link for http://www.apple.com
        When I ask the system for the redirection list
        Then the system returns a list containing the short links