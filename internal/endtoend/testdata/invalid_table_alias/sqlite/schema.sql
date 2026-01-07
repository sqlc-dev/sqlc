-- https://github.com/sqlc-dev/sqlc/issues/437
CREATE TABLE authors (
  id   INTEGER PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  bio  text
);
