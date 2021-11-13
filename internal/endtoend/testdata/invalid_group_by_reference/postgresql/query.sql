CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL,
  bio  text
);

-- name: ListAuthors :many
SELECT   *
FROM     authors
GROUP BY invalid_reference;
