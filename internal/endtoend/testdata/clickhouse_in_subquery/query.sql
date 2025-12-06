-- name: GetUsersInDepartments :many
SELECT id, name FROM users WHERE department_id IN (SELECT id FROM departments WHERE name = ?);

-- name: GetUsersNotInDepartments :many
SELECT id, name FROM users WHERE department_id NOT IN (SELECT id FROM departments WHERE name IN (?, ?));

-- name: GetUsersWithIds :many
SELECT id, name FROM users WHERE id IN (?, ?, ?);
