SELECT id FROM users WHERE id IN (sqlc.slice(ids)) AND name = @name;
