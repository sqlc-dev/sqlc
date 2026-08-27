-- name: GetAuthor :one
select  id,name , bio from   authors
where id =  ? limit 1;

# hash comment
-- name: ListAuthors :many
SELECT id, name, bio FROM authors
ORDER BY name DESC;

-- name: PickyQuery :many
SELECT id, -- the primary key
       name,
       -- computed downstream
       bio
FROM authors
-- soft-deleted rows are filtered
WHERE bio IS NOT NULL
  AND id > ?;

-- name: InlineBlock :many
SELECT /* inline note */ id, name FROM authors ORDER BY name;

-- name: CountSigils :one
SELECT count(*) FROM authors
WHERE id <> ? AND name <> @user_name; # session variable stays

-- name: CastUnsigned :one
SELECT CAST(id AS UNSIGNED) FROM authors LIMIT 1;

-- name: FirstTwin :one
SELECT id FROM authors LIMIT 1;

-- name: SecondTwin :one
SELECT id FROM authors LIMIT 1;

-- name: CreateAuthor :execresult
insert into authors (
  name, bio
) values (
  ?, ?
);
