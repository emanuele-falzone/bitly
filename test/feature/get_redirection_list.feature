# language: en

Feature: Get redirection list
    As an administrator
    I want to know the list of short link
    So that I can see if there are useless links

    Scenario: The list contains the expected links
        Given that I got a short link for http://www.google.com
        And that I got a short link for http://www.apple.com
        When I ask the system for the redirection list
        Then the system returns a list containing the short links