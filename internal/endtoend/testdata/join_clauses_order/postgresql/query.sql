CREATE TABLE a (
  id   BIGSERIAL PRIMARY KEY,
  a TEXT NOT NULL
);

CREATE TABLE b (
  id BIGSERIAL PRIMARY KEY,
  b TEXT NOT NULL,
  a_id BIGINT NOT NULL REFERENCES a (id)
);

CREATE TABLE c (
  id BIGSERIAL PRIMARY KEY,
  c TEXT NOT NULL,
  a_id BIGINT NOT NULL REFERENCES a (id)
);

CREATE TABLE d (
  id BIGSERIAL PRIMARY KEY,
  d TEXT NOT NULL,
  a_id BIGINT NOT NULL REFERENCES a (id)
);

CREATE TABLE e (
  id BIGSERIAL PRIMARY KEY,
  e TEXT NOT NULL,
  a_id BIGINT NOT NULL REFERENCES a (id)
);

-- name: TestLeftInner :many
SELECT a.a, b.b, c.c
FROM a
LEFT JOIN b ON b.a_id = a.id
INNER JOIN c ON c.a_id = a.id;

-- name: TestInnerLeft :many
SELECT a.a, b.b, c.c
FROM a
INNER JOIN b ON b.a_id = a.id
LEFT JOIN c ON c.a_id = a.id;

-- name: TestLeftInnerLeftInner :many
SELECT a.a, b.b, c.c, d.d, e.e
FROM a
LEFT JOIN b ON b.a_id = a.id
INNER JOIN c ON c.a_id = a.id
LEFT JOIN d ON d.a_id = a.id
INNER JOIN e ON e.a_id = a.id;

-- name: TestInnerLeftInnerLeft :many
SELECT a.a, b.b, c.c, d.d, e.e
FROM a
INNER JOIN b ON b.a_id = a.id
LEFT JOIN c ON c.a_id = a.id
INNER JOIN d ON d.a_id = a.id
LEFT JOIN e ON e.a_id = a.id;
