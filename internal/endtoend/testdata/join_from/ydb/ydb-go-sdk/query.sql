-- name: MultiFrom :many
SELECT email FROM bar, foo WHERE login = $login;
