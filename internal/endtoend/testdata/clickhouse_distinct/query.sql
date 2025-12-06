-- name: DistinctDepartments :many
SELECT DISTINCT department FROM users;

-- name: DistinctMultipleColumns :many
SELECT DISTINCT department, salary FROM users ORDER BY department;
