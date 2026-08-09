-- name: SkipVetAll :exec
-- @sqlc-vet-disable always-fail no-exec
SELECT true FROM users WHERE id = $1;
