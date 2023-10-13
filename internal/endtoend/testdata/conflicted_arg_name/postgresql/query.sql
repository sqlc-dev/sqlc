-- name: Time2ByTime :one
SELECT time2 FROM foo WHERE time=$1;

-- name: Uuid2ByUuid :one
SELECT uuid2 FROM foo WHERE uuid=$1;
