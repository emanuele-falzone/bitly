# language: en

Feature: Navigate to a short URL
    As a visitor
    I want to know the location of a short link
    So that I can visit the specified resource

    Scenario: Navigate to a previouly created short URL
        Given that I got a short link for http://www.google.com
        When I navigate to the short link
        Then the system redirects me to http://www.google.com
    
    Scenario: Navigate to a short URL that does not exist
        Given that I got a short link that does not exist
        When I navigate to the short link
        Then the system signals that the short link does not exist