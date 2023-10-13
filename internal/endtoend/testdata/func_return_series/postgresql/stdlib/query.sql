/* name: GenerateSeries :many */
SELECT ($1::int) + i
FROM generate_series(0, $2::int) AS i
LIMIT 1;
