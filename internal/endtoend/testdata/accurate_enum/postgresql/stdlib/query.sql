-- name: ListTasks :many
SELECT * FROM tasks;

-- name: GetTasksByStatus :many
SELECT * FROM tasks WHERE status = $1;

-- name: CreateTask :one
INSERT INTO tasks (title, status) VALUES ($1, $2) RETURNING *;
