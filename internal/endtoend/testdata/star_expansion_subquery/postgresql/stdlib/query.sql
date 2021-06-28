CREATE TABLE foo (a text, b text);

-- name: StarExpansionSubquery :many
SELECT * FROM foo WHERE EXISTS (SELECT * FROM foo);
