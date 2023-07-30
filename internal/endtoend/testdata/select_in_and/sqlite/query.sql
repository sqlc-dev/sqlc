-- Example queries for sqlc
CREATE TABLE authors (
  id   integer PRIMARY KEY,
  name text      NOT NULL,
  age  integer
);

CREATE TABLE translators (
  id   integer PRIMARY KEY,
  name text      NOT NULL,
  age  integer
);

CREATE TABLE books (
  id   integer PRIMARY KEY,
  author text      NOT NULL,
  translator text      NOT NULL,
  year  integer
);

-- name: DeleteAuthor :exec
DELETE FROM
  books AS b
WHERE
  b.author NOT IN (
    SELECT
      a.name
    FROM
      authors a
    WHERE
      a.age >= ?
  )
  AND b.translator NOT IN (
    SELECT
      t.name
    FROM
      translators t
    WHERE
      t.age >= ?
  )
  AND b.year <= ?;