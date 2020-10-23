/* name: ListCities :many */
SELECT *
FROM city
ORDER BY name;

/* name: GetCity :one */
SELECT *
FROM city
WHERE slug = ?;

/* name: CreateCity :execresult */
INSERT INTO city (
    name,
    slug,
    data
) VALUES (
    ?,
    ?,
    ?
); 

/* name: UpdateCityName :exec */
UPDATE city
SET name = ?
WHERE slug = ?;
