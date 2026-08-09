SELECT a.*, b.name FROM a JOIN b ON b.id = a.b_id WHERE a.id = $1;
