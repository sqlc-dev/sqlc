CREATE TABLE things (
  id integer NOT NULL,
  title varchar,
  price decimal(10, 2) NOT NULL,
  tags varchar[],
  grid integer[][] NOT NULL,
  point struct(x integer, y integer),
  attrs map(varchar, integer),
  fixed integer[3]
);
