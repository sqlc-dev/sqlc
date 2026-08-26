-- name: GetAuthor :one
select  id,name , bio from   authors
where id =  ? limit 1;

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

-- name: ListAuthors :many
SELECT id, name, bio FROM authors
ORDER BY name;

-- name: GetEvent :one
SELECT "EventName", "order"
FROM "Events" WHERE id = ? LIMIT 1;

-- name: CountSigils :one
SELECT count(*) FROM authors
WHERE id <> @at_param AND id <> :colon_param AND id <> $dollar_param;

-- name: TopAuthors :many
WITH ranked AS (
  SELECT id, name FROM authors
)
SELECT * FROM ranked;

-- name: AddAuthor :exec
INSERT INTO authors (name, bio)
VALUES (?, ?);

-- name: MakeScratch :exec
CREATE TABLE scratch (
  id INTEGER NOT NULL,
  label TEXT
);
