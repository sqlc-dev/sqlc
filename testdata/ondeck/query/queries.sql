-- name: ListCities
-- result: many
SELECT *
FROM city
ORDER BY name;

-- name: GetCity
-- result: one
SELECT *
FROM city
WHERE slug = $1;

-- name: ListVenues
-- result: many
SELECT *
FROM venue
WHERE city = $1
ORDER BY name;
