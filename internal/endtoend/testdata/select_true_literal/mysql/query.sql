-- name: SelectTrue :one
SELECT true;

-- name: SelectFalse :one
SELECT false;

-- name: SelectTrueWithAlias :one
SELECT true AS is_active;

-- name: SelectMultipleBooleans :one
SELECT true AS col_a, false AS col_b, true AS col_c;
