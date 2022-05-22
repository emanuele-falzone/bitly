# language: en

Feature: Inspect how many times a short link has been visited
    As an administrator
    I want to know how many time a given short link has been visited
    So that I can see if the link I shared has been used or not

    Scenario: The short link exists
        Given that I got a short link for http://www.google.com
        And that the link has been visited 5 times
        When I ask the system how many times the link has been visited
        Then the system says the link has been visited 5 times
    
    Scenario: The short link does not exist
        Given that I got a short link that does not exist
        When I ask the system how many times the link has been visited
        Then the system signals that the short link does not exist