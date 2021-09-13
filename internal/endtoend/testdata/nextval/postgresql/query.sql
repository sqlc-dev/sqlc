-- Simple table
CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY
);

-- name: GetNextID :one
SELECT pk, pk FROM
 (SELECT nextval('authors_id_seq') as pk) AS alias;

