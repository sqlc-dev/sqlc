/* name: GetUsers :many */
SELECT *
FROM users_func()
WHERE first_name != '';
