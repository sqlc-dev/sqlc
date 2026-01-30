-- name: ColdCities :many
SELECT city
FROM weather
GROUP BY city
HAVING Max(temp_lo) < $max_temp;
