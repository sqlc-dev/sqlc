-- name: GetByPublicID :one
SELECT *
FROM locations l
WHERE l.public_id = ?
AND EXISTS (
    SELECT 1
    FROM projects p
    WHERE p.id = location.project_id
);
