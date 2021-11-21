CREATE TABLE foo (id text not null);

CREATE FUNCTION get_id(input text) RETURNS text AS $$ SELECT 'bar' $$ LANGUAGE sql;

-- name: ListFoos :one
SELECT id FROM foo WHERE id = get_id(CASE WHEN $1 = 100 THEN $2 ELSE 'baz' END);
