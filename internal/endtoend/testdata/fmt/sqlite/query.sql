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

-- name: MakeMeasurements :exec
CREATE TABLE IF NOT EXISTS measurements (
  id INTEGER PRIMARY KEY,
  label VARYING  CHARACTER(120),
  ratio DECIMAL(10,5),
  note
);

-- name: KeepAutoinc :exec
CREATE TABLE   counters (id INTEGER PRIMARY KEY AUTOINCREMENT, hits INTEGER DEFAULT 0);

-- name: CastLabel :one
SELECT CAST(bio AS VARYING CHARACTER(120)) FROM authors LIMIT 1;

-- name: MakeSearchIndex :exec
CREATE  VIRTUAL  TABLE IF NOT EXISTS notes_fts USING fts5(body, tag UNINDEXED, tokenize = 'porter');

-- name: MakeRecipeIndex :exec
CREATE VIRTUAL TABLE recipes_fts USING fts5(
  name,
  ingredients
);

-- name: LoginName :one
SELECT COALESCE(bio, '') AS login FROM authors WHERE id = ? LIMIT 1;

-- name: SpelledJoins :many
SELECT a.id
FROM authors AS a
INNER JOIN authors AS b ON a.id = b.id
LEFT OUTER JOIN authors AS c ON a.id = c.id
WHERE a.id != ? AND b.id == ? AND c.id <> ?;

-- name: CommaJoin :one
SELECT count(*) FROM authors AS a, authors AS b WHERE a.id = b.id;

-- name: PlannerHint :many
SELECT a.id FROM authors AS a CROSS JOIN authors AS b;

-- name: NumberedParams :many
SELECT id FROM authors WHERE name = ?2 AND bio = ?1;

-- name: NewestAndOldest :many
SELECT id FROM authors
UNION
SELECT id FROM authors;

-- name: QuickUnion :many
SELECT id FROM authors UNION ALL SELECT id FROM authors;

-- name: EnableForeignKeys :exec
PRAGMA foreign_keys = 1;

-- name: UpsertAuthor :exec
INSERT INTO authors (id, name)
VALUES (?, ?)
ON CONFLICT (id) DO UPDATE SET
  name = excluded.name
WHERE excluded.name <> '';

-- name: QuickUpsert :exec
INSERT INTO authors (id, name) VALUES (?, ?) ON CONFLICT (id) DO UPDATE SET name = excluded.name;
