-- name: GreaterThan :many
SELECT COUNT(*) > 0 FROM bar;

-- name: LessThan :many
SELECT COUNT(*) < 0 FROM bar;

-- name: GreaterThanOrEqual :many
SELECT COUNT(*) >= 0 FROM bar;

-- name: LessThanOrEqual :many
SELECT COUNT(*) <= 0 FROM bar;

-- name: NotEqual :many
SELECT COUNT(*) != 0 FROM bar;

-- name: AlsoNotEqual :many
SELECT COUNT(*) <> 0 FROM bar;

-- name: Equal :many
SELECT COUNT(*) = 0 FROM bar;


