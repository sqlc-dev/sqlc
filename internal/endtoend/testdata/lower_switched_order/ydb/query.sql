-- name: LowerSwitchedOrder :many
SELECT bar FROM foo WHERE bar = $bar AND bat = Unicode::ToLower($bat_lower);
