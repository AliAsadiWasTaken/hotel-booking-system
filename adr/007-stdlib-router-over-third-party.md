# ADR-0007

## Title

stdlib net/http router over third-party routing libraries.

## Status

Accepted

## Context

Go has a large ecosystem of HTTP routing libraries including Gin, Chi,
Echo, and Fiber. These libraries provide features such as middleware
chaining, route groups, parameter extraction helpers, and response
utilities that are not available in the standard library.

Prior to Go 1.22, the stdlib ServeMux had significant limitations:
it could not match HTTP methods or extract path parameters, making
third-party routers a practical necessity for REST APIs.

Go 1.22 introduced method-based routing and path parameter extraction
directly in net/http/ServeMux, closing the gap for the majority of
REST API use cases.

## Decision

Use the stdlib net/http ServeMux introduced in Go 1.22.

Routes are registered with method and path patterns:

    mux.HandleFunc("GET /hotels/{id}", handler.GetByID)

Path parameters are extracted with:

    r.PathValue("id")

## Consequences

Advantages

- Zero additional dependencies for routing.
- No framework-specific abstractions to learn or debug around.
- Route patterns are readable and self-documenting.
- The standard library is stable, well-documented, and maintained
  by the Go team.

Disadvantages

- No built-in middleware chaining — middleware must be composed
  manually using handler wrapping.
- No route groups — common prefixes must be repeated per route
  registration.
- Missing convenience features such as automatic OPTIONS handling
  and built-in body size limits. These must be implemented explicitly
  where needed.
