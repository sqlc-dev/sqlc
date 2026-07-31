SELECT * FROM users WHERE a = ?1 AND b = ?2 AND c IN (/*SLICE:ids*/?);
