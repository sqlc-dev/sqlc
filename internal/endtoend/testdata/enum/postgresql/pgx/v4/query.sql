-- name: GetAll :many
SELECT * FROM users;

-- name: NewUser :exec
INSERT INTO users (
    first_name,
    last_name,
    age,
    shoe_size,
    shirt_size
) VALUES
($1, $2, $3, $4, $5);

-- name: UpdateSizes :exec
UPDATE users
SET shoe_size = $2, shirt_size = $3
WHERE id = $1;

-- name: DeleteBySize :exec
DELETE FROM users
WHERE shoe_size = $1 AND shirt_size = $2;