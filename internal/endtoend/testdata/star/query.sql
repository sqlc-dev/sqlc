CREATE TABLE bar (bid serial not null);
CREATE TABLE foo (fid serial not null);

-- name: Star :many
SELECT * FROM bar, foo;
