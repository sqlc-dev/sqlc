CREATE TABLE x (
  a   text,
  b   text
);

CREATE TABLE y (
  a   text,
  b   text
);

-- name: DeleteXWithY :exec
DELETE x FROM x INNER JOIN y ON y.a = x.a WHERE x.b = y.b;