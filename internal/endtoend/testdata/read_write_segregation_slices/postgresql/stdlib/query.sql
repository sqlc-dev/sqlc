/* name: ListUsersByIDs :many */
SELECT
    id, first_name
FROM users
WHERE id = ANY($1::int[]);

/* name: GetUserTags :one */
SELECT tags FROM users WHERE id = $1;

/* name: SetUserTags :exec */
UPDATE users SET tags = $1 WHERE id = $2;
