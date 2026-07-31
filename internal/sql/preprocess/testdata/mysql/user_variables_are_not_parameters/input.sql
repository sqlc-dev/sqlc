SELECT @row := @row + 1 FROM users WHERE id = sqlc.arg(id);
