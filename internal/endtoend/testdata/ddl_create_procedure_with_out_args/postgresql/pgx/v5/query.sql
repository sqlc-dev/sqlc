-- name: CallInsertData :one
CALL insert_data($1, $2, null);

-- name: CallInsertDataNoArgs :exec
CALL insert_data(1, 2, null);

-- name: CallInsertDataNamed :one
CALL insert_data(b => $1, a => $2, c => null);

-- name: CallInsertDataSqlcArgs :exec
CALL insert_data(sqlc.arg('foo'), sqlc.arg('bar'), sqlc.arg('с'));
