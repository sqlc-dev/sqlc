-- name: InsertCode :one
WITH cc AS (
            INSERT INTO td3.codes(created_by, updated_by, code, hash, is_private)
            VALUES ($1, $1, $2, $3, false)
            RETURNING hash
)
INSERT INTO td3.test_codes(created_by, updated_by, test_id, code_hash)
VALUES(
            $1, $1, $4, (select hash from cc)
)
RETURNING *;
