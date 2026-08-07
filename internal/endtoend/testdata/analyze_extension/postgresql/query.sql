-- name: GetUserByEmail :one
SELECT id, email FROM users WHERE email = $1;

-- name: MatchingEmails :many
SELECT id FROM users WHERE email = name;

-- name: RankUsers :many
SELECT id, similarity(name, $1) AS score FROM users;
