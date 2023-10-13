-- name: UpdateCode :one
WITH cc AS (
            UPDATE td3.codes
            SET
                created_by = $1,
                updated_by  = $1,
                code = $2,
                hash = $3,
                is_private = false
            RETURNING hash
)
UPDATE td3.test_codes
SET
    created_by = $1,
    updated_by  = $1,
    test_id = $4,
    code_hash = cc.hash
    FROM cc
RETURNING *;
