INSERT INTO t (a) VALUES (sqlc.arg(val))
ON DUPLICATE KEY UPDATE a = sqlc.arg(new_val);
