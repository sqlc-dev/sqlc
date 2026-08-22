-- name: GetAuthor :one
select  id,name , bio from   authors
where id =  ? limit 1;

# hash comment
-- name: ListAuthors :many
SELECT id, name, bio FROM authors
ORDER BY name DESC;

-- name: CreateAuthor :execresult
insert into authors (
  name, bio
) values (
  ?, ?
);
