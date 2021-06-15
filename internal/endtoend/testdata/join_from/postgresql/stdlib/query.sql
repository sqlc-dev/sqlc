CREATE TABLE foo (email text not null);
CREATE TABLE bar (login text not null);

-- name: MultiFrom :many
SELECT email FROM bar, foo WHERE login = $1;
