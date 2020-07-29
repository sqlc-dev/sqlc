CREATE FUNCTION f(x TIMESTAMPTZ) RETURNS void AS '' LANGUAGE sql;
CREATE FUNCTION f(x timestamp with time zone) RETURNS void AS '' LANGUAGE sql;

-- name: F :one
SELECT f(1);
