-- name: GetGroups :many
SELECT
    rg.groupId,
    rg.groupName
FROM 
    routergroup rg
WHERE
    rg.groupName LIKE '%' || COALESCE($groupName, rg.groupName) || '%' AND
    rg.groupId = COALESCE($groupId, rg.groupId);

