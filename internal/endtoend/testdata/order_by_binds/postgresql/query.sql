CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL,
  bio  text
);

-- name: ListAuthorsColumnSort :many
SELECT  * FROM authors
WHERE   id > sqlc.arg(min_id) 
ORDER   BY CASE WHEN sqlc.arg(sort_column) = 'name' THEN name END;

-- name: ListAuthorsNameSort :many
SELECT  * FROM authors
WHERE   id > sqlc.arg(min_id)
ORDER   BY name ASC;
