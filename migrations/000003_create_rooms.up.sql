CREATE TABLE rooms (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT        NOT NULL,
    hotel_id    UUID        NOT NULL REFERENCES hotels(id) ON DELETE CASCADE,
    bed_count   INTEGER     NOT NULL,
    capacity    INTEGER     NOT NULL,
    quantity    INTEGER     NOT NULL,
    price_per_night NUMERIC NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
