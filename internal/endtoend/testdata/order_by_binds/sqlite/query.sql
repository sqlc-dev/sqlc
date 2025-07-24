-- name: ListAuthorsColumnSort :many
SELECT  * FROM authors
WHERE   id > sqlc.arg(min_id) 
ORDER   BY CASE WHEN sqlc.arg(sort_column) = 'name' THEN name END;

-- name: ListAuthorsColumnSortDirection :many
SELECT * FROM authors
WHERE id > ?
ORDER BY
    CASE
        WHEN @order_by = 'asc' THEN name
    END ASC,
    CASE
        WHEN @order_by = 'desc' OR @order_by IS NULL THEN name
    END DESC;

-- name: ListAuthorsColumnSortFnWtihArg :many
SELECT  * FROM authors
ORDER   BY id % ?;

-- name: ListAuthorsNameSort :many
SELECT  * FROM authors
WHERE   id > sqlc.arg(min_id)
ORDER   BY name ASC;
