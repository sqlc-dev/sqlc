-- name: EncodeDigest :one
SELECT encode(digest($1, 'sha1'), 'hex');

