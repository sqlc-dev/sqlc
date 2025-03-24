CREATE TABLE authors (
  id   INT  PRIMARY KEY,
  name text NOT NULL,
  bio  text NOT NULL
);

ALTER TABLE authors ADD COLUMN explanation text NOT NULL DEFAULT ('');
