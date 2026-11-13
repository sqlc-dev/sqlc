-- name: NestedCTEScope :one
WITH items AS (
  SELECT 1::integer AS value
),
derived AS (
  SELECT value FROM (
    WITH items AS (
      SELECT 'inner'::text AS value
    )
    SELECT value FROM items
  ) AS nested
)
SELECT items.value AS outer_value, derived.value AS nested_value
FROM items CROSS JOIN derived;
