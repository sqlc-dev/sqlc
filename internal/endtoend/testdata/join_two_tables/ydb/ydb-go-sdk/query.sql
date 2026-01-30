-- name: TwoJoins :many
SELECT foo.*
FROM foo
JOIN bar ON bar.id = foo.bar_id
JOIN baz ON baz.id = foo.baz_id;
