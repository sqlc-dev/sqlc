SELECT sqlc.embed(a), b.name FROM a JOIN b ON b.id = a.b_id WHERE a.id = sqlc.arg(id);
