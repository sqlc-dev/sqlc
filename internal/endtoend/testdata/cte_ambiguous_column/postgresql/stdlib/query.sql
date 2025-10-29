-- name: GetWorkspacesJoinTasks :many
WITH wtask AS (
    SELECT
        workspaces.*,
        tasks.id IS NOT NULL::boolean AS has_task
    FROM workspaces
    LEFT JOIN tasks ON tasks.workspace_id = workspaces.id
)
SELECT *
FROM wtask
ORDER BY CASE WHEN owner_id = @owner_id THEN 0 ELSE 1 END;

-- name: GetWorkspacesJoinTasksRenameColumn :many
WITH wtask AS (
    SELECT
        workspaces.owner_id AS w_owner_id,
        workspaces.*,
        tasks.id IS NOT NULL::boolean AS has_task
    FROM workspaces
    LEFT JOIN tasks ON tasks.workspace_id = workspaces.id
)
SELECT *
FROM wtask
ORDER BY CASE WHEN w_owner_id = @owner_id THEN 0 ELSE 1 END;

-- name: GetWorkspacesSubQueryTasks :many
WITH wfiltered AS (
    SELECT workspaces.*
    FROM workspaces
    WHERE EXISTS (
        SELECT 1
        FROM tasks
        WHERE tasks.workspace_id = workspaces.id
    )
)
SELECT *
FROM wfiltered
ORDER BY CASE WHEN owner_id = @owner_id THEN 0 ELSE 1 END;
