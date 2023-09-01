CREATE FUNCTION f$n() RETURNS integer 
    AS $$ SELECT 1 $$ LANGUAGE SQL;

-- name: Fn :one
SELECT f$n();
