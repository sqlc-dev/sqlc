CREATE TABLE bar (id int not null primary key autoincrement);

-- name: BarExists :one
SELECT
    EXISTS (
        SELECT
            1
        FROM
            bar
        where
            id = ?
    );
