# ADR-0006

## Title

Database transactions for booking creation and cancellation.

## Status

Accepted

## Context

Booking creation involves multiple dependent database operations: reading
the room to check its quantity, counting overlapping confirmed bookings,
and inserting the new booking. These operations must behave as a single
unit.

Without a transaction, two concurrent requests targeting the same room
and dates can both read the same availability state, both conclude there
is capacity, and both insert a booking — exceeding the room's quantity.
This is a classic read-modify-write race condition.

Cancellation has a similar problem: two concurrent cancellation requests
can both read status = "confirmed", both pass the guard check, and both
attempt to cancel an already-cancelled booking.

## Decision

Wrap booking creation and cancellation in database transactions.

For booking creation, a SELECT ... FOR UPDATE is issued against the room
row before the availability check. This acquires a row-level lock that
blocks any other transaction attempting to lock the same row, serializing
concurrent requests for the same room. The lock is released when the
transaction commits or rolls back.

For cancellation, the booking row is read and the status guard is
evaluated inside a transaction. The transaction ensures the read and the
subsequent update are atomic.

The DBTX interface, satisfied by both *pgxpool.Pool and pgx.Tx, allows
repositories to operate inside a transaction without any changes to their
implementation. The service layer begins the transaction, constructs
transaction-scoped repositories by passing the tx, and commits or rolls
back based on the outcome.

## Consequences

Advantages

- Overbooking is prevented at the database level regardless of
  concurrent request volume or application server count.
- Duplicate cancellations are rejected atomically rather than relying
  on application-level state that can become stale between requests.
- Repository code is unchanged — the transaction boundary is owned
  entirely by the service layer, which is the appropriate place for
  business operation boundaries.
- defer tx.Rollback ensures the transaction is always cleaned up even
  if an unexpected error or panic occurs mid-operation.

Disadvantages

- Row-level locks on the room during booking creation increase contention
  for popular rooms under high concurrency. Requests queue behind the
  lock holder rather than running in parallel.
- Longer transactions hold locks longer. Any slow operation inside the
  transaction (network call, heavy computation) would delay other
  requests targeting the same room.
