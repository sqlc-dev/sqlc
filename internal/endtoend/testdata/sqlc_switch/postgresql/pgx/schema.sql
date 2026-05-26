CREATE TABLE authors (
    id         BIGSERIAL PRIMARY KEY,
    name       text NOT NULL,
    created_at timestamptz NOT NULL
);
