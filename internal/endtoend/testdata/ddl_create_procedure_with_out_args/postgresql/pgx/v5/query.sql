-- name: CallInsertData :one
CALL insert_data($1, $2, null, null, null, null, null, null, null, null, null, null, null);

-- name: CallInsertDataNoArgs :one
CALL insert_data(1, 2, null, null, null, null, null, null, null, null, null, null, null);

-- name: CallInsertDataNamed :one
CALL insert_data(
        b => $1,
        a => $2,
        c => null,
        i => null,
        j => null,
        k => null,
        d => null,
        h => null,
        e => null,
        m => null,
        f => null,
        g => null,
        l => null
     );

-- name: CallInsertDataSqlcArgs :one
CALL insert_data(sqlc.arg('foo'), sqlc.arg('bar'), null, null, null, null, null, null, null, null, null, null, null);
