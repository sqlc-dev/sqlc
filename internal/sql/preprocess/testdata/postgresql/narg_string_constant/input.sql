SELECT * FROM users WHERE name = COALESCE(sqlc.narg('name'), name);
