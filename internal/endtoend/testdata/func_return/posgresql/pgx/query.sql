
/* name: GetUsers :many */
SELECT *
FROM users_func()
WHERE first_name != '';

/* name: GenerateSeries :many */
SELECT ($1::inet) + i
FROM generate_series(0, $2::int) AS i
LIMIT 1;