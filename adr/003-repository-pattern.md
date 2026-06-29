# ADR-0003

## Title

Repository pattern to separate data access from business logic.

## Status

Accepted

## Context

The application needs to read and write data to Postgres. The question
is where SQL queries should live and how tightly coupled business logic
should be to the database.

Without a repository layer, SQL queries would be written directly inside
service or handler functions, mixing infrastructure concerns with
business rules.

## Decision

Introduce a repository layer as the only place in the codebase that
executes SQL queries.

Services receive and return domain types. They have no knowledge of
SQL, connection pools, or Postgres-specific behaviour.

Repositories translate between domain types and database rows.

## Consequences

Advantages

- Business logic is decoupled from the database. The service layer can
  be tested independently by substituting a fake repository.
- If the storage backend changes, only the repository changes. Services
  and handlers are untouched.
- SQL is co-located by feature. All hotel queries live in
  hotel/repository.go, making them easy to find and audit.
- The boundary makes it obvious where a query belongs. There is never
  a question of whether to put SQL in the handler or the service.

Disadvantages

- Adds a layer of indirection for simple read operations where no
  business logic exists.
- Requires explicit wiring of repository dependencies into services
  and services into handlers, since there is no dependency injection
  framework.
