-- Example queries for sqlc
CREATE TABLE authors (
  id   INTEGER PRIMARY KEY,
  name text      NOT NULL
);

CREATE TABLE books (
  id INTEGER PRIMARY KEY,
  title text NOT NULL
);

