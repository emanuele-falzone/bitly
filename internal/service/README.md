# Service

The URL Shortening service relies on different services.

## Key generator
The key service generator abstracts the computation of a given `key` starting from the `location` value.
THere can be different implementation. For example the key can be completely random, or computed based on the location itself using hashing functions such ad SHA1 or MD5.

## Event store
This service is in change of permanently store the dispatched events in order to enable further processing.

## Event logger
This service is in change of logging the events so as to enable the developer to easily inspect the application behavior.