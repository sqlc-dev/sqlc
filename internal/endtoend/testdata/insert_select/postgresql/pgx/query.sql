CREATE TABLE bar (name text not null, ready bool not null);
CREATE TABLE foo (name text not null, meta text not null);

-- name: InsertSelect :exec
INSERT INTO foo (name, meta)
SELECT name, $1
FROM bar WHERE ready = $2;
