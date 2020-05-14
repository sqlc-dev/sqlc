-- name: MakeIntervalSecs :one
SELECT make_interval(secs => $1);

-- name: MakeIntervalDays :one
SELECT make_interval(days => $1::int);

-- name: MakeIntervalMonths :one
SELECT make_interval(months => sqlc.arg('months')::int);
