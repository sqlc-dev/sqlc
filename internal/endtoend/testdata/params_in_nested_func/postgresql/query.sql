-- name: GetGroups :many
SELECT
    rg.groupId,
    rg.groupName
FROM 
    RouterGroup rg
WHERE
    rgr.depth = 1 AND
    rg.groupName LIKE CONCAT('%', COALESCE(sqlc.narg('groupName')::text, rg.groupName), '%') AND
    rg.groupId = COALESCE(sqlc.narg('groupId'), rg.groupId);
