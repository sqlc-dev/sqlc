-- name: Active :macro
SELECT * FROM authors WHERE deleted_at IS NULL;

-- name: ListBooks :many
WITH a AS sqlc.macro.Active()
SELECT sqlc.embed(b), a.name FROM a JOIN books b ON b.author_id = a.id WHERE b.id > sqlc.arg(after);
