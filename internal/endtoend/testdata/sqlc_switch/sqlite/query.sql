-- name: FindAuthors :many
SELECT * FROM authors
WHERE sqlc.switch(@filter,
    sqlc.when('named', 'name IS NOT NULL'),
    sqlc.else(         '1 = 1'));
