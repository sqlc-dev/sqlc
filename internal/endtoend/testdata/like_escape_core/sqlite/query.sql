-- name: SearchName :many
SELECT id FROM t WHERE name LIKE ? ESCAPE ?;

-- name: SearchNameNamed :many
SELECT id FROM t WHERE name LIKE @pattern ESCAPE @esc;

-- name: SearchNotLike :many
SELECT id FROM t WHERE name NOT LIKE ? ESCAPE ?;

-- name: SearchEscapeOnly :many
SELECT id FROM t WHERE name LIKE 'a!%%' ESCAPE ?;

-- name: SearchLiteralEscape :many
SELECT id FROM t WHERE name LIKE ? ESCAPE '!';

-- name: SearchThenWhere :many
SELECT id FROM t WHERE name LIKE ? ESCAPE ? AND score > ?;
