-- https://www.postgresql.org/docs/current/sql-syntax-calling-funcs.html
CREATE FUNCTION concat_lower_or_upper(a text, b text, uppercase boolean DEFAULT false)
RETURNS text
AS
$$
 SELECT CASE
        WHEN $3 THEN UPPER($1 || ' ' || $2)
        ELSE LOWER($1 || ' ' || $2)
        END;
$$
LANGUAGE SQL IMMUTABLE STRICT;

-- name: PositionalNotation :one
SELECT concat_lower_or_upper('Hello', 'World', true);

-- name: PositionalNoDefaault :one
SELECT concat_lower_or_upper('Hello', 'World');

-- name: NamedNotation :one
SELECT concat_lower_or_upper(a => 'Hello', b => 'World');

-- name: NamedAnyOrder :one
SELECT concat_lower_or_upper(a => 'Hello', b => 'World', uppercase => true);

-- name: NamedOtherOrder :one
SELECT concat_lower_or_upper(a => 'Hello', uppercase => true, b => 'World');

-- name: MixedNotation :one
SELECT concat_lower_or_upper('Hello', 'World', uppercase => true);
