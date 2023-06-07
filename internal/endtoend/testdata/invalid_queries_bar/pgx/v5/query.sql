CREATE TABLE foo (bar text not null, baz text not null);
INSERT INTO foo (bar, baz) VALUES ($1);
INSERT INTO foo (bar) VALUES ($1, $2);
