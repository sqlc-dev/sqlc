-- name: Subtree :many
WITH tree AS (
    SELECT id, parent_id, name, 0 AS depth FROM nodes WHERE id = @root
    UNION ALL
    SELECT n.id, n.parent_id, n.name, t.depth + 1
    FROM nodes n
    JOIN tree t ON n.parent_id = t.id
)
SELECT id, parent_id, name, depth FROM tree
OPTION (MAXRECURSION 100);

-- name: Numbers :many
WITH seq(i) AS (
    SELECT 1
    UNION ALL
    SELECT i + 1 FROM seq WHERE i < @max
)
SELECT i FROM seq;
