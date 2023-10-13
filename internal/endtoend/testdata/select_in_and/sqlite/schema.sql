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

