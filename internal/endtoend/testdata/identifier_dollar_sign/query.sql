CREATE FUNCTION f$n() RETURNS integer AS 'SELECT 1';

-- name: Fn :one
SELECT f$n();
