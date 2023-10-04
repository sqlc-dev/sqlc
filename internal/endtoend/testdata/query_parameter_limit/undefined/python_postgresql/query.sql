-- name: DeleteBarByID :execrows
DELETE FROM bar WHERE id = $1;

-- name: DeleteBarByIDAndName :execrows
DELETE FROM bar
WHERE id = $1
AND name1 = $2
AND name2 = $3
AND name3 = $4
;
