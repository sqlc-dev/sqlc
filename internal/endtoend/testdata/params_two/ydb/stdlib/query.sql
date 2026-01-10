-- name: FooByAandB :many
SELECT a, b FROM foo 
WHERE a = $a and b = $b;
