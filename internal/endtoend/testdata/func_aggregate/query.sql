-- Example queries for sqlc
CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL,
  bio  text
);

-- name: Percentile :one
select percentile_disc(0.5) within group (order by authors.name)
from authors;
