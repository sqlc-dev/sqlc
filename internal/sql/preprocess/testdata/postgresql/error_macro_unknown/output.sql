                     
         

-- name: Bad :one
SELECT * FROM sqlc.macro.Unknown(sqlc.arg(id)) u;

-- name: Good :one
SELECT * FROM (SELECT 1) k WHERE id = $1;

-- The name of a macro is matched exactly.
-- name: WrongCase :one
SELECT * FROM sqlc.macro.known() k;
