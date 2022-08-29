CREATE TABLE weather (
  city     text    NOT NULL,
  temp_lo  integer NOT NULL
);

-- name: ColdCities :many
SELECT city
FROM weather
GROUP BY city
HAVING max(temp_lo) < ?;
