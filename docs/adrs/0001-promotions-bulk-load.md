# Promotions Bulk Loading With Redis CLI

## Context and Problem Statement

How to load millions or even more entries into Redis in an efficient way?

## Considered Options

### Redis pipeline using Go Redis client

Using a regular Redis client for Go, like `go-redis`, resulted in too much round trips for each operation with Redis,
which degraded the performance.

### In memory solution

An in memory solution would require an in house implementation, requiring time and prone to errors.
Given the fact that there are well established, lightweight solutions, such as Redis, the in memory solution was
discarded.

### Other database

Other database solutions like Postgres, MySQL were considered but quickly discarded considering the nature of the data,
simple, temporary, requires high availability. Managing this data with an in memory solution showed to be better.

## Decision Outcome

Chosen option was `bulk load with Redis CLI`.
**redis-cli --pipe** solution allows us to stream each of the promotions to Redis in a much more efficient way,
considering that the command will only return a summary of the streamed data at the end.

# References

- [Bulk Loading](https://redis.io/docs/reference/patterns/bulk-loading/)
