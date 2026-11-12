SELECT id FROM users WHERE id IN (/*SLICE:ids*/@ids) AND name = @name;
