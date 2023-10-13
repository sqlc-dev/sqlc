-- https://github.com/sqlc-dev/sqlc/issues/1590
CREATE TABLE authors (
  name text NOT NULL,
  deleted_at datetime NOT NULL,
  created_at datetime NOT NULL,
  updated_at datetime NOT NULL
);

CREATE TABLE books (
  is_amazing tinyint(1) NOT NULL,
  deleted_at datetime NOT NULL,
  created_at datetime NOT NULL,
  updated_at datetime NOT NULL
);

