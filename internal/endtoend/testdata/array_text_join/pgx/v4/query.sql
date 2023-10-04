-- name: JoinTextArray :many
SELECT bar.info
FROM foo
JOIN bar ON foo.bar = bar.id;
