CREATE TABLE foo (bar text);

-- name: DoubleDash :one
SELECT * FROM foo WHERE bar = $1;

/* name: SlashStar :one */
SELECT * FROM foo WHERE bar = $1;
