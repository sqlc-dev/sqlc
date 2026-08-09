SELECT id FROM users WHERE id = ANY(sqlc.slice(ids));
