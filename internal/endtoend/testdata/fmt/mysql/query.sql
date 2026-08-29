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

-- name: CasePreserved :many
SELECT ID, Name FROM Authors WHERE Name <> '' ORDER BY Name;

-- name: LiteralsSurvive :one
SELECT true AS t, false AS f, 1.50 AS score, CASE WHEN name = '' THEN NULL ELSE name END AS n FROM authors LIMIT 1;

-- name: DistinctNames :many
SELECT DISTINCT name FROM authors;

-- name: UnionOrdered :many
SELECT name AS foo FROM authors
UNION
SELECT bio AS foo FROM authors
ORDER BY foo;

-- name: ConcatNames :many
SELECT group_concat(DISTINCT name ORDER BY name DESC SEPARATOR ' ') FROM authors GROUP BY bio;

-- name: Hinted :one
SELECT /*+ MAX_EXECUTION_TIME(1000) */ id FROM authors LIMIT 1;

-- name: UpdateWithJoin :exec
UPDATE authors AS a JOIN authors AS p ON p.id = a.id
SET a.name = ?
WHERE p.bio IS NOT NULL;

-- name: ShowWarnings :many
SHOW WARNINGS;

CREATE TABLE scores (points decimal(10, 5), views bigint unsigned NOT NULL);
