-- name: GetGroups :many
SELECT
    rg.groupId,
    rg.groupName
FROM 
    RouterGroup rg
WHERE
    rg.groupName LIKE CONCAT('%', COALESCE(sqlc.narg('groupName'), rg.groupName), '%') AND
    rg.groupId = COALESCE(sqlc.narg('groupId'), rg.groupId);
