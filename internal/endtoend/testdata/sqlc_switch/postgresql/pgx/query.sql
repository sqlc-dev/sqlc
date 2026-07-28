-- name: ListAuthors :many
SELECT * FROM authors
WHERE name = $1
ORDER BY sqlc.switch(@sort,
    sqlc.when('name_asc', 'authors.name ASC'),
    sqlc.when('recent',   'authors.created_at DESC, authors.id DESC'),
    sqlc.else(            'authors.id ASC'));
