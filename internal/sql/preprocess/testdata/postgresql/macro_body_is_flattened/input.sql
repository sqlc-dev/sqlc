-- The body lands on one line, without its comments, so the query it is pasted
-- into keeps its lines and gains no annotations.
-- name: Recent :macro
-- @param n int
SELECT id, -- the key
       'not -- a comment' AS note, /* left
out */ $tag$ /* nor this */ $tag$ AS tagged
FROM   events
ORDER  BY id DESC
LIMIT  sqlc.arg(n);

-- name: ListRecent :many
SELECT *
FROM sqlc.macro.Recent(10) r
WHERE r.id > sqlc.arg(after);
