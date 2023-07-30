CREATE TABLE bar (id integer not null primary key autoincrement);

-- name: BarNotExists :one
SELECT 
    NOT EXISTS (
        SELECT
            1
        FROM
            bar
        WHERE
            id = ?
    );

