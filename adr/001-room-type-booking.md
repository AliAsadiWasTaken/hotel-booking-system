# ADR-0001

## Title

Bookings reference room types instead of physical rooms.

## Status

Accepted

## Context

The booking platform needs to determine whether reservations should target
individual hotel rooms or room types.

Hotels typically assign physical rooms during check-in, not during online
booking.

## Decision

Bookings will reference room types.

Each room type has a quantity that represents the inventory available.

Availability is determined by comparing overlapping bookings against the room
type quantity.

## Consequences

Advantages

- Simpler data model.
- Simpler booking algorithm.
- Matches how many hotel booking platforms work.

Disadvantages

- Cannot assign physical room numbers.
- A hotel PMS would be needed for room assignment.