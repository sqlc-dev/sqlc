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
