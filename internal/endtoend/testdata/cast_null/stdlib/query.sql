-- name: ListNullable :many
SELECT
  NULL::text as a,
  NULL::integer as b,
  NULL::bigint as c,
  NULL::time as d
FROM foo;
