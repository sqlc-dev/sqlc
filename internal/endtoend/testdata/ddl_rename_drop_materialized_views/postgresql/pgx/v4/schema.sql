CREATE TABLE authors (
    id   BIGSERIAL PRIMARY KEY,
    name text      NOT NULL,
    bio  text
);

CREATE MATERIALIZED VIEW authors_mv AS (
  SELECT * FROM authors
);

CREATE MATERIALIZED VIEW authors_mv_new AS (
  SELECT * FROM authors
);

ALTER MATERIALIZED VIEW authors_mv RENAME TO authors_mv_old;
ALTER MATERIALIZED VIEW authors_mv_new RENAME TO authors_mv;

DROP MATERIALIZED VIEW authors_mv_old;
