# Promotions Bulk Loading With Redis CLI

## Context and Problem Statement

How to load millions of even more entries into Redis in an efficient way?

## Considered Options

- Redis pipeline using Go Redis client
- Bulk load using Redis CLI

## Decision Outcome

Chosen option was `bulk load with Redis CLI`.
Even having to write a file containing all Redis commands, that becomes bigger than the original csv, using the Redis
cli to send to data to the server showed to be much faster than Go client.

# References

- [Bulk Loading](https://redis.io/docs/reference/patterns/bulk-loading/)
