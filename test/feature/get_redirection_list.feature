# language: en

Feature: Get redirection list

    Scenario: The list contains the expected links
        Given that I got a short link for http://www.google.com
        And that I got a short link for http://www.apple.com
        When I ask the system for the redirection list
        Then the system returns a list containing the short links