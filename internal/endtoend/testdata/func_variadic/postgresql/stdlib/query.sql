-- name: UpdateJ :exec
UPDATE
    test
SET
    j = jsonb_build_object($1::text, $2::text, $3::text, $4::text)
WHERE
    id = $5;