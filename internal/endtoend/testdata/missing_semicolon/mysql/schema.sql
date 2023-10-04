-- https://github.com/sqlc-dev/sqlc/issues/1198
CREATE TABLE authors (
  id   INT PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  bio  text
);

