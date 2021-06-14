CREATE TABLE foo (id text not null, bar text not null);
CREATE TABLE bar (id text not null, info text[] not null);

-- name: JoinTextArray :many
SELECT bar.info
FROM foo
JOIN bar ON foo.bar = bar.id;
