CREATE TABLE authors (
  id        INT PRIMARY KEY,
  name      TEXT NOT NULL,
  parent_id INT -- nullable
);

-- name: AllAuthors :many
SELECT  a.id,
        a.name,
        p.id as alias_id,
        p.name as alias_name
FROM    authors a
        LEFT JOIN authors p USING (parent_id);
