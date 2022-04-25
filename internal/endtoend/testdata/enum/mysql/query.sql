/* name: GetAll :many */
SELECT * FROM users;

-- name: NewUser :exec
INSERT INTO users (
    id,
    first_name,
    last_name,
    age,
    shoe_size,
    shirt_size
) VALUES
(?, ?, ?, ?, ?, ?);

-- name: UpdateSizes :exec
UPDATE users
SET shoe_size = ?, shirt_size = ?
WHERE id = ?;

-- name: DeleteBySize :exec
DELETE FROM users
WHERE shoe_size = ? AND shirt_size = ?;
