-- name: FindByAddress :one
SELECT * FROM "user" WHERE "metadata"->>'address1' = ?1 LIMIT 1;