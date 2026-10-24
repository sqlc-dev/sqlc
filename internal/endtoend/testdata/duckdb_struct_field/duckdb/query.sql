-- name: Bare :many
SELECT info.label FROM items;

-- name: Qualified :many
SELECT i.info.label AS label, i.info.score AS score FROM items AS i;

-- name: Filtered :many
SELECT id FROM items AS i WHERE i.info.score > $1;

-- name: CaseInsensitive :many
SELECT info.LABEL AS label FROM items;
