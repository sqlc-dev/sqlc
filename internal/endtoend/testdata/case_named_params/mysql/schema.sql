-- https://github.com/sqlc-dev/sqlc/issues/1195

CREATE TABLE authors (
  id   BIGINT PRIMARY KEY,
  username VARCHAR(10) NULL,
  email VARCHAR(10) NULL,
  name TEXT  NOT NULL,
  bio  TEXT,
  UNIQUE KEY idx_username (username),
  UNIQUE KEY ids_email (email)
);

