-- https://github.com/sqlc-dev/sqlc/issues/1195

CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  username TEXT NULL,
  email TEXT NULL,
  name TEXT  NOT NULL,
  bio  TEXT
);

