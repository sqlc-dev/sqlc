-- name: GetAuthor :one
select  id,name , bio from   authors
where id =  ? limit 1;

-- name: ListAuthors :many
SELECT id, name, bio FROM authors
ORDER BY name;
