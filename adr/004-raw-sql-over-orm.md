# ADR-0004

## Title

Raw SQL with pgx over an ORM.

## Status

Accepted

## Context

Go has several popular database libraries. At one end, ORMs like GORM
generate SQL automatically from struct definitions and method chains.
At the other end, drivers like pgx execute raw SQL strings directly.

The project needs database access that is predictable, debuggable, and
able to express Postgres-specific queries such as RETURNING, date
overlap logic, and aggregations.

## Decision

Use pgx directly with raw SQL strings. No ORM.

All queries are written explicitly. The exact SQL sent to Postgres is
always visible in the source code.

## Consequences

Advantages

- Every query is readable and auditable without understanding ORM
  internals or query generation rules.
- Postgres-specific features (RETURNING, gen_random_uuid, DATE types,
  aggregations, window functions) can be used without ORM workarounds.
- Performance problems are diagnosed at the SQL level, not by
  inspecting generated output from an abstraction layer.
- Engineers can talk about exactly what queries the system runs, which
  is valuable in technical interviews and code reviews.

Disadvantages

- More boilerplate per query compared to an ORM. Each query requires
  explicit column lists in both SELECT and Scan.
- No automatic schema-to-struct mapping. Column order in SELECT must
  match argument order in Scan.
- Schema changes require manually updating queries in the repository.
