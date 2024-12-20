-- name: GetFullNames :many
SELECT
 full_name
FROM
  (
    SELECT
  	 concat(first_name, ' ', last_name) as full_name
    FROM
      customers
  ) subquery
WHERE
  full_name IN (sqlc.slice ("full_names"));

-- name: GetFullNames2 :many
WITH subquery AS (
    SELECT
  	 concat(first_name, ' ', last_name) as full_name
    FROM
      customers
  )
SELECT
 full_name
FROM
  subquery
WHERE
  full_name IN (sqlc.slice ("full_names"));