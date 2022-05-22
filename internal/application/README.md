# Application

The application is composed by several commands and queries that allows for abstraction domain objects and services.

The whole business logic relies under the application package.
In this way we can test user requirements directly interacting with application code and relying on in memory repositories.
Such an attitude ensure that no business logic leaked in the adapter layer.