-- name: InsertVector :exec
INSERT INTO foo(embedding) VALUES (STRING_TO_VECTOR('[0.1, 0.2, 0.3, 0.4]'));

-- name: SelectVector :many
SELECT id FROM foo
ORDER BY DISTANCE(STRING_TO_VECTOR('[1.2, 3.4, 5.6]'), embedding, 'L2_squared')
LIMIT 10;
