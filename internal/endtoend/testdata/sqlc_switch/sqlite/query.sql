-- name: FindAuthors :many
SELECT * FROM authors
WHERE sqlc.switch(@filter,
    sqlc.when('named', 'name IS NOT NULL'),
    sqlc.else(         '1 = 1'));

-- name: ListAuthors :many
SELECT id, name, created_at FROM authors
ORDER BY sqlc.switch(@sort,
    sqlc.when('name_asc', 'authors.name ASC'),
    sqlc.when('recent',   'authors.created_at DESC, authors.id DESC'),
    sqlc.else(            'authors.id ASC'));
