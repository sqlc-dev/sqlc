-- name: AllAuthors :many
SELECT
  a.id,
  a.name,
  p.id AS alias_id,
  p.name AS alias_name
FROM authors AS a
LEFT JOIN authors AS p ON a.parent_id = p.id;
