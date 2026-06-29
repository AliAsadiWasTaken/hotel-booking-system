# ADR-0008

## Title

Struct-tag validation with go-playground/validator.

## Status

Accepted

## Context

HTTP handlers receive untrusted input that must be validated before
being passed to the service layer. Validation can be implemented in
two ways: manually with explicit if-statements per field, or
declaratively using struct tags processed by a validation library.

The service layer's CreateInput structs intentionally contain no
validation logic — their job is to carry data between layers, not
enforce rules. Validation belongs at the boundary where untrusted
data enters the system, which is the handler layer.

## Decision

Use github.com/go-playground/validator with struct tags on request
body types defined inside handler functions.

Validation rules are expressed as struct tags:

    Name  string `validate:"required"`
    Email string `validate:"required,email"`
    Price float64 `validate:"required,gt=0"`

A shared Validate function in the api package runs validation and
returns the first failing field as a human-readable message.

## Consequences

Advantages

- Validation rules are co-located with the fields they apply to,
  making them easy to read and modify.
- A single Validate call covers all fields, avoiding scattered
  if-statements per field.
- go-playground/validator is the de facto standard validation library
  in the Go ecosystem with broad tag support.

Disadvantages

- Only the first validation error is returned per request. Clients
  receive one error at a time rather than a full list of all failing
  fields. This can be improved later if the API requires it.
- Struct tags are strings — validation rules are not checked by the
  compiler and typos fail silently at runtime.
