CREATE TABLE foo (
  bar text,
  bam name,
  baz name not null
);

-- name: ListFoo :many
SELECT * FROM foo;
