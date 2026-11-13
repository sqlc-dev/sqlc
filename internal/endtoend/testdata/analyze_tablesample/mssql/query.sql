-- name: Sampled :many
SELECT id FROM t TABLESAMPLE (@pct PERCENT) WHERE id > @min;

-- name: SampledRows :many
SELECT id FROM t TABLESAMPLE SYSTEM (@n ROWS) REPEATABLE (@seed);
