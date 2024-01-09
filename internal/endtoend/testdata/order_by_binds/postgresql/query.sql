-- name: ListAuthorsColumnSort :many
SELECT  * FROM authors
WHERE   id > sqlc.arg(min_id) 
ORDER   BY CASE WHEN sqlc.arg(sort_column) = 'name' THEN name END;

-- name: ListAuthorsColumnSortFnWtihArg :many
SELECT  * FROM authors
ORDER   BY MOD(id, $1);

-- name: ListAuthorsNameSort :many
SELECT  * FROM authors
WHERE   id > sqlc.arg(min_id)
ORDER   BY name ASC;
