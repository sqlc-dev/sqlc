-- name: CallInsertData :exec
CALL insert_data($1, $2);

-- name: CallInsertDataNoArgs :exec
CALL insert_data(1, 2);

-- name: CallInsertDataNamed :exec
CALL insert_data(b => $1, a => $2);

-- name: CallInsertDataSqlcArgs :exec
CALL insert_data(sqlc.arg('foo'), sqlc.arg('bar'));
