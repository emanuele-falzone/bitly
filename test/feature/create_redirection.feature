# language: en

Feature: Create a short link for given link
    As an administrator
    I want to create a short link
    So that I can easily share a resource with other people

    Scenario: The given link is valid
        When I command the system to shorten the link http://www.google.com
        Then the system returns a short link

    Scenario: The given link is malformed
        When I command the system to shorten the link google.com
        Then the system signals that the link is not valid