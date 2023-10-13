-- name: CallInsertData :exec
CALL insert_data(?, ?);

-- name: CallInsertDataNoArgs :exec
CALL insert_data(1, 2);

-- name: CallInsertDataSqlcArgs :exec
CALL insert_data(sqlc.arg('foo'), sqlc.arg('bar'));
