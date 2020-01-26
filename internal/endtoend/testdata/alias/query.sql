CREATE TABLE bar (id serial not null);
CREATE TABLE foo (id serial not null, bar serial references bar(id));

-- name: Alias :exec
DELETE FROM foo f USING bar b
WHERE f.bar = b.id AND b.id = $1;
