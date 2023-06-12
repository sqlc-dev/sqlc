-- name: GetValues :many
SELECT id, index::bigint, value::text
FROM array_values AS x, unnest(values) WITH ORDINALITY AS y (value, index);
