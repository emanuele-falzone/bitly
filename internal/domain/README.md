# Domain

In this document I will try to describe the reasoning behind such a domain modeling.

I think that a URL Shortening service relies on a simple entity, here named `Redirection`.
A redirection is composed by a `Location`, i.e. the URL the user want to shorten, and a `Key`,
i.e. the string the system use as an abbreviation for the location.

We can also formally define events that are related to the life cycle of a redirection:
- `created`: the redirection has been successfully created and permanently stored inside the repository.
- `read`: the redirection has been requested by a user by specifying the redirection key.
- `deleted`: the redirection has been permanently deleted from the repository.

Such events can be `dispatched` to a series of `listeners` and possibly analyzed in order to extract information, such as the number or times a redirection has been read or the average TTL of redirections.
