CREATE TABLE x (
  a   text,
  b   text
);

CREATE TABLE y (
  a   text,
  b   text
);

-- name: UpdateXWithY :exec
UPDATE x INNER JOIN y ON y.a = x.a SET x.b = y.b;
