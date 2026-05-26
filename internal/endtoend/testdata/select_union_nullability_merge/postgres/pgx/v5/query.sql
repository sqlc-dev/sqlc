-- UNION ALL — column nullability is OR-merged across legs:
-- left leg's `id`/`label` are NOT NULL, right leg's are nullable → result nullable.
-- `rank` is NOT NULL on both sides → result NOT NULL.
-- name: UnionAllMixed :many
SELECT id, label, rank FROM not_null_leg
UNION ALL
SELECT id, label, rank FROM nullable_leg;

-- UNION (set-distinct) — same nullability merge rule applies.
-- name: UnionMixed :many
SELECT id, label, rank FROM not_null_leg
UNION
SELECT id, label, rank FROM nullable_leg;

-- Chained UNION ALL — middle and right legs introduce nullability.
-- name: UnionAllChained :many
SELECT id, label, rank FROM not_null_leg
UNION ALL
SELECT id, label, rank FROM nullable_leg
UNION ALL
SELECT id, label, rank FROM not_null_leg;

-- Both legs NOT NULL — result stays NOT NULL.
-- name: UnionAllBothNotNull :many
SELECT id, label, rank FROM not_null_leg
UNION ALL
SELECT id, label, rank FROM not_null_leg;
