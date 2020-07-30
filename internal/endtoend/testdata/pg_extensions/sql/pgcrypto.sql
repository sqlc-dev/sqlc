CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- name: EncodeDigest :one
SELECT encode(digest($1, 'sha1'), 'hex');

