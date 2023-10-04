-- name: LowerSwitchedOrder :many
SELECT bar FROM foo WHERE bar = $1 AND bat = LOWER($2);
