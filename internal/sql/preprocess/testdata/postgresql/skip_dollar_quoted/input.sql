CREATE FUNCTION f() RETURNS int AS $$ SELECT sqlc.arg(nope); $$ LANGUAGE sql;
