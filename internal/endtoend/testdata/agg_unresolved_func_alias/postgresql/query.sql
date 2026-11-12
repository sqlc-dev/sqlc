-- name: UnresolvedFuncAlias :exec
SELECT ARRAY_AGG(u.*) FROM some_unknown_func() AS u;
