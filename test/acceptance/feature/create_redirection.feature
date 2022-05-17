# language: en

Feature: Create a short link for given link

    Scenario: The given link is valid
        When I command the system to shorten the link http://www.google.com
        Then the system returns a short link

    Scenario: The given link is malformed
        When I command the system to shorten the link google.com
        Then the system signals that the link is not valid