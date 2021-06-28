CREATE TABLE foo (
    other     text NOT NULL,
    total     bigint NOT NULL,
    tags      text[] NOT NULL,
    byte_seq  bytea NOT NULL,
    retyped   text NOT NULL,
    langs     text[]
);
