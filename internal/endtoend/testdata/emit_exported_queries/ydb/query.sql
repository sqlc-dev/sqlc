-- name: UpdateBarID :exec
UPDATE bar SET id = $new_id WHERE id = $old_id;
