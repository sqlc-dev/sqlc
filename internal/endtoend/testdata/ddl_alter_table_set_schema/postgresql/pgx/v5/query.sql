-- name: GetFooBar :exec
SELECT * FROM foo.bar;

-- name: UpdateFooBar :exec
UPDATE foo.bar SET name = $1;
