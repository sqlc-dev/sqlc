CREATE TABLE authors (
  id   integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  name text    NOT NULL,
  bio  text,
  UNIQUE(name)
);
