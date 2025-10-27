-- name: GetData :many
SELECT *
FROM my_table
WHERE (CAST($allow_invalid AS Bool) OR NOT invalid);



