-- name: SessionInfo :one
SELECT @@ROWCOUNT AS row_count,
    @@IDENTITY AS last_identity,
    @@SPID AS session_id,
    @@VERSION AS version,
    @@TRANCOUNT AS tran_count,
    @@ERROR AS error,
    @@SERVERNAME AS server_name,
    @@DATEFIRST AS date_first,
    @@DBTS AS dbts;

-- name: TouchInTransaction :exec
UPDATE readings SET note = @note WHERE id = @id AND @@TRANCOUNT > 0;
