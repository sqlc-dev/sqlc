-- name: GreaterThan :many
SELECT count(*) > 0 FROM bar;

-- name: LessThan :many
SELECT count(*) < 0 FROM bar;

-- name: GreaterThanOrEqual :many
SELECT count(*) >= 0 FROM bar;

-- name: LessThanOrEqual :many
SELECT count(*) <= 0 FROM bar;

-- name: NotEqual :many
SELECT count(*) != 0 FROM bar;

-- name: AlsoNotEqual :many
SELECT count(*) <> 0 FROM bar;

-- name: Equal :many
SELECT count(*) = 0 FROM bar;






