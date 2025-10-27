-- name: ListAuthorsColumnSort :many
SELECT * FROM authors
WHERE id > $min_id 
ORDER BY CASE WHEN $sort_column = 'name' THEN name END;

-- name: ListAuthorsColumnSortFnWtihArg :many
SELECT * FROM authors
ORDER BY Math::mod(id, $mod_arg);

-- name: ListAuthorsNameSort :many
SELECT * FROM authors
WHERE id > $min_id
ORDER BY name ASC;

