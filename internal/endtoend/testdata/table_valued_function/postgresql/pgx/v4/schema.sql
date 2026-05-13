CREATE TABLE authors (
  id   int PRIMARY KEY,
  name text      NOT NULL
);

CREATE FUNCTION fauthors() returns table(
  id   int,
  name text
 ) stable as $$
 	select * from authors
    $$ language sql;
