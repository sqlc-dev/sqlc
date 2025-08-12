/* name: UpdateSet :exec */
UPDATE foo SET name = ? WHERE slug = ?;

/* name: UpdateSetQuoted :exec */
UPDATE "foo" SET "name" = ? WHERE "slug" = ?;
