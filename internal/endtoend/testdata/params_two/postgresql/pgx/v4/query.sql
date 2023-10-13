-- name: FooByAandB :many
SELECT a, b FROM foo 
WHERE a = $1 and b = $2;
