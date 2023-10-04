-- name: FindAccount :one
SELECT
    a.*,
    ua.name
    -- other fields
FROM
    accounts a
    INNER JOIN users_accounts ua ON a.id = ua.id2
WHERE
    a.id = @account_id;
