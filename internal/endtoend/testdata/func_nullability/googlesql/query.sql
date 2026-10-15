-- name: Aggregates :one
SELECT COUNT(*) AS total, SUM(amount) AS total_amount, AVG(id) AS mean_id, MAX(name) AS last_name,
  ARRAY_AGG(id) AS ids, STRING_AGG(name, ',') AS names
FROM events;

-- name: Propagates :many
SELECT UPPER(name) AS upper_name, UPPER(nickname) AS upper_nickname, LENGTH(nickname) AS nickname_length
FROM events;

-- name: MayBeNull :many
SELECT JSON_VALUE(details, '$.kind') AS kind, REGEXP_EXTRACT(name, 'a+') AS run,
  SAFE_DIVIDE(id, 0) AS ratio, SAFE.SUBSTR(name, 1, 3) AS prefix, ARRAY_FIRST(scores) AS first_score
FROM events;

-- name: Conditionals :many
SELECT COALESCE(nickname, name) AS display_name, IFNULL(ratio, 0) AS ratio_or_zero,
  NULLIF(name, '') AS non_empty, IF(id > 10, name, nickname) AS chosen, GREATEST(id, 1) AS at_least_one
FROM events;

-- name: ParamsFromFunctions :many
SELECT id FROM events
WHERE STARTS_WITH(name, @prefix) AND TIMESTAMP_TRUNC(happened_at, DAY) = @day
  AND EXTRACT(DATE FROM happened_at) = @on AND ARRAY_INCLUDES(scores, @score);
