# sqlc Documentation

> And lo, the Great One looked down upon the people and proclaimed:
> "SQL is actually pretty great"

sqlc generates **fully type-safe idiomatic Go code** from SQL. Here's how it
works:

1. You write SQL queries
2. You run sqlc to generate Go code that presents type-safe interfaces to those
   queries
3. You write application code that calls the methods sqlc generated

Seriously, it's that easy. You don't have to write any boilerplate SQL querying
code ever again.

## Getting started

- [Installing sqlc](overview/install.md)
- [Getting started with MySQL](tutorials/getting-started-mysql.md)
- [Getting started with PostgreSQL](tutorials/getting-started-postgresql.md)
- [Getting started with SQLite](tutorials/getting-started-sqlite.md)
