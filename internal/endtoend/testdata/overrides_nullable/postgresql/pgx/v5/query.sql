CREATE TABLE foo (
  bar text,
  bam jsonb,
  baz jsonb not null
);

-- name: ListFoo :many
SELECT * FROM foo;
