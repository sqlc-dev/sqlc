-- name: GetFirst :many
SELECT * FROM first_view;

-- name: GetSecond :many
SELECT * FROM second_view WHERE val2 = $1;
