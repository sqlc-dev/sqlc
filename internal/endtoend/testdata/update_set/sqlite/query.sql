CREATE TABLE foo (name text not null, slug text not null);

/* name: UpdateSet :exec */
UPDATE foo SET name = ? WHERE slug = ?;
