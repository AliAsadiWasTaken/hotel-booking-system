# ADR-0002

## Title

Feature-first package structure over package-by-layer.

## Status

Accepted

## Context

Go projects can be organized in two primary ways.

Package-by-layer groups code by its technical role.

```
internal/
  repositories/
  services/
  handlers/
```

Package-by-feature groups code by the domain concept it belongs to.

```
internal/
  hotel/
  room/
  booking/
  user/
```

The project needs a structure that scales as features grow and that
makes eventual microservice extraction straightforward.

## Decision

Organize packages by feature.

Each domain package owns its model, repository, service, and handler.

## Consequences

Advantages

- All code related to a feature lives in one place. Working on hotels
  means working inside the hotel package only.
- Microservice extraction becomes moving a package, not untangling
  cross-cutting dependencies across layer directories.
- Mirrors how NestJS modules work by convention, making the mental
  model familiar for engineers coming from that background.
- Import paths communicate intent: hotel.Repository is clearly the
  hotel domain's data access, not a generic repository.

Disadvantages

- Shared infrastructure (database connection, logger) still needs to
  be passed across packages, which requires deliberate wiring in main.
- New contributors may expect a layers-based structure if they come
  from frameworks that enforce it.
