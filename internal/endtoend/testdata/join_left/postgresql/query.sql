--- https://github.com/kyleconroy/sqlc/issues/604
CREATE TABLE users (
  user_id    INT PRIMARY KEY,
  city_id    INT -- nullable
);
CREATE TABLE cities (
  city_id    INT PRIMARY KEY,
  mayor_id   INT NOT NULL
);
CREATE TABLE mayors (
  mayor_id   INT PRIMARY KEY,
  full_name  TEXT NOT NULL
);

-- name: GetMayors :many
SELECT
    user_id,
    mayors.full_name
FROM users
LEFT JOIN cities USING (city_id)
INNER JOIN mayors USING (mayor_id);

-- name: GetMayorsOptional :many
SELECT
    user_id,
    mayors.full_name
FROM users
LEFT JOIN cities USING (city_id)
LEFT JOIN mayors USING (mayor_id);

-- name: GetMayorsOptionalInnerSelect :many
SELECT t1.user_id, t2.full_name
FROM (
    SELECT user_id FROM users WHERE users.city_id = $1 LIMIT 1 OFFSET 0
) AS t1
LEFT JOIN cities on users.city_id = cities.city_id
LEFT JOIN (
    SELECT mayors.mayor_id, mayors.full_name
    FROM mayors where mayors.mayor_id = $2
) AS t2 on t2.mayor_id = cities.mayor_id;
