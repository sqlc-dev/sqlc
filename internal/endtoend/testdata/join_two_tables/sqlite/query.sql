CREATE TABLE foo (bar_id integer not null, baz_id integer not null);
CREATE TABLE bar (id integer not null);
CREATE TABLE baz (id integer not null);

-- name: TwoJoins :many
SELECT foo.*
FROM foo
JOIN bar ON bar.id = bar_id
JOIN baz ON baz.id = baz_id;
