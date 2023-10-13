CREATE FUNCTION foo(bar TEXT, baz TEXT='baz') RETURNS bool AS $$ SELECT true $$ LANGUAGE sql;
