-- name: AllAuthors :many
SELECT  a.id,
        a.name,
        p.id as alias_id,
        p.name as alias_name
FROM    authors a
        LEFT JOIN authors p
            ON (authors.parent_id = p.id);
