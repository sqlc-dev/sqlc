-- name: GetSomeDeletedNotOk :many
DELETE FROM a
USING b
WHERE a.b_id_fk = b.b_id
RETURNING b.b_id; -- column "b_id" does not exist
