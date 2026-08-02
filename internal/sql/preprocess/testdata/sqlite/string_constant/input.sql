SELECT * FROM users WHERE a = sqlc.arg('x') AND b = sqlc.narg('y') AND c IN (sqlc.slice('ids'));
