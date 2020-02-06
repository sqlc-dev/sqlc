CREATE TABLE foo (bar text not null, baz text not null);
INSERT INTO foo (bar, baz) VALUES ($1);
INSERT INTO foo (bar) VALUES ($1, $2);

-- stderr
-- # package querytest
-- query.sql:2:1: INSERT has more target columns than expressions
-- query.sql:3:1: INSERT has more expressions than target columns
