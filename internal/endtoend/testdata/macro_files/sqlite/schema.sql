CREATE TABLE authors (
  id         integer PRIMARY KEY,
  name       text NOT NULL,
  deleted_at datetime
);

CREATE TABLE books (
  id         integer PRIMARY KEY,
  author_id  integer NOT NULL,
  title      text NOT NULL,
  deleted_at datetime
);
