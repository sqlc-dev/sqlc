-- name: UpdateJ :exec
UPDATE
    test
SET
    j = jsonb_build_object($1, $2, $3, $4)
WHERE
    id = $5;