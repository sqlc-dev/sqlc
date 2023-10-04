-- name: ColdCities :many
SELECT city
FROM weather
GROUP BY city
HAVING max(temp_lo) < $1;
