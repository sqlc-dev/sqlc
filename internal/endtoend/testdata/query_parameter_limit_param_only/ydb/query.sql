-- name: MarkNoticeDone :exec
UPDATE notice
SET status='done', notice_at=$notice_at
WHERE id=$id;

-- name: CreateNotice :exec
INSERT INTO notice (cnt, created_at)
VALUES ($cnt, $created_at);
