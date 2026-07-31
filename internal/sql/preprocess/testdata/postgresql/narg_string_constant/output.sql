SELECT * FROM users WHERE name = COALESCE($1, name);
