-- name: MarkNoticeDone :exec
UPDATE notice
SET status='done', notice_at=$1
WHERE id=$2;

-- name: CreateNotice :exec
INSERT INTO notice (cnt, created_at)
VALUES ($1, $2);