-- name: Ranking :many
SELECT
  id,
  LAG(price) OVER (ORDER BY id) AS previous_price,
  LEAD(title, 1, 'none') OVER (ORDER BY id) AS next_title,
  FIRST_VALUE(subtitle) OVER (ORDER BY id) AS first_subtitle,
  LAST_VALUE(qty) OVER (ORDER BY id) AS last_qty
FROM books;

-- name: Nulls :many
SELECT
  ISNULL(rating, 0) AS rating,
  ISNULL(discount, 0) AS discount,
  NULLIF(qty, 0) AS qty,
  CHOOSE(1, title, subtitle) AS chosen,
  GREATEST(price, discount) AS greatest_price,
  LEAST(title, subtitle) AS least_title
FROM books;

-- name: Strings :many
SELECT
  REPLACE(title, 'a', 'b') AS replaced,
  REPLACE(note, 'a', 'b') AS replaced_n,
  REVERSE(title) AS reversed,
  CHARINDEX('a', title) AS position,
  PATINDEX('%a%', note) AS pattern_position,
  FORMAT(price, 'N2') AS formatted,
  QUOTENAME(name) AS quoted,
  STR(weight, 10, 2) AS str_weight
FROM books;

-- name: Aggregates :one
SELECT
  STRING_AGG(title, ',') AS titles,
  STRING_AGG(note, ',') AS notes,
  STDEV(qty) AS qty_stdev
FROM books;

-- name: Math :many
SELECT
  ABS(qty) AS abs_qty,
  ABS(price) AS abs_price,
  ABS(cost) AS abs_cost,
  CEILING(weight) AS ceiling_weight,
  FLOOR(rating) AS floor_rating,
  ROUND(discount, 1) AS rounded,
  POWER(weight, 2) AS squared,
  SQRT(qty) AS root,
  SIGN(discount) AS sign_discount
FROM books;

-- name: Dates :many
SELECT
  YEAR(published) AS year,
  MONTH(updated) AS month,
  DAY(published) AS day,
  EOMONTH(published) AS month_end,
  DATEFROMPARTS(rating, 1, 1) AS start_date,
  DATETIME2FROMPARTS(rating, 1, 1, 0, 0, 0, 0, 3) AS start_time,
  SWITCHOFFSET(offset_at, '+01:00') AS switched
FROM books;

-- name: Json :many
SELECT
  ISJSON(doc) AS is_json,
  JSON_VALUE(doc, '$.a') AS json_value,
  JSON_QUERY(doc, '$.a') AS json_query,
  JSON_MODIFY(doc, '$.a', 1) AS json_modified,
  HASHBYTES('SHA2_256', title) AS hash,
  CHECKSUM(title) AS checksum
FROM books;
