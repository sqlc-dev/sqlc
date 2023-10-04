CREATE TABLE authors (
    id   BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    bio  TEXT
);

ALTER TABLE authors ADD COLUMN gender INTEGER NULL;

CREATE MATERIALIZED VIEW authors_names as SELECT name from authors;

