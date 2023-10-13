-- name: TwoJoins :many
SELECT foo.*
FROM foo
JOIN bar ON bar.id = bar_id
JOIN baz ON baz.id = baz_id;
