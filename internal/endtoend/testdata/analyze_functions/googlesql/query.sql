-- name: Strings :many
SELECT LENGTH(name) AS len, UPPER(name) AS upper_name, CONCAT(name, '!') AS shout,
  SPLIT(name, ',') AS parts, STARTS_WITH(name, 'a') AS starts, TO_HEX(SHA256(name)) AS digest,
  FARM_FINGERPRINT(name) AS fingerprint, REPLACE(name, 'a', 'b') AS replaced, TO_BASE64(payload) AS b64
FROM events;

-- name: Numbers :many
SELECT ABS(id) AS abs_id, id / 2 AS half, amount * id AS scaled, ratio + id AS shifted,
  ROUND(ratio, 2) AS rounded, MOD(id, 7) AS bucket, DIV(id, 2) AS halved, COSINE_DISTANCE(embedding, embedding) AS distance
FROM events;

-- name: Times :many
SELECT DATE_ADD(happened_on, INTERVAL 1 DAY) AS next_day, DATE_DIFF(happened_on, happened_on, DAY) AS days,
  TIMESTAMP_TRUNC(happened_at, HOUR) AS hour, EXTRACT(YEAR FROM happened_at) AS year,
  EXTRACT(DATE FROM happened_at) AS day, UNIX_SECONDS(happened_at) AS epoch,
  FORMAT_TIMESTAMP('%Y', happened_at) AS formatted
FROM events;

-- name: Collections :many
SELECT ARRAY_LENGTH(scores) AS n, ARRAY_REVERSE(scores) AS reversed, GENERATE_ARRAY(1, id) AS series,
  TO_JSON(id) AS as_json, JSON_TYPE(details) AS json_kind
FROM events;

-- name: Counts :one
SELECT COUNT(*) AS total, COUNT(name) AS named, COUNTIF(ratio > 0.5) AS high FROM events;
