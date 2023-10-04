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

