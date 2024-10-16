-- name: RunVetAll :exec
SELECT true;

-- name: SkipVetAll :exec
-- @sqlc-vet-disable
SELECT true;

-- name: SkipVetSingleLine :exec
-- @sqlc-vet-disable always-fail no-exec
SELECT true;

-- name: SkipVetMultiLine :exec
-- @sqlc-vet-disable always-fail
-- @sqlc-vet-disable no-exec
SELECT true;

-- name: SkipVet_always_fail :exec
-- @sqlc-vet-disable always-fail
SELECT true;

-- name: SkipVet_no_exec :exec
-- @sqlc-vet-disable no-exec
SELECT true;

-- name: SkipVetInvalidRule :exec
-- @sqlc-vet-disable always-fail
-- @sqlc-vet-disable block-delete
-- @sqlc-vet-disable no-exec
SELECT true;