-- name: CreateTodoNullPlaceholder :exec
CALL create_todo(sqlc.arg(task)::text, null);

-- name: CreateTodoTypedNullPlaceholder :exec
CALL create_todo(sqlc.arg(task)::text, NULL::int);

-- name: CreateTodoNamedOut :exec
CALL create_todo(sqlc.arg(task)::text, sqlc.arg(id)::int);
