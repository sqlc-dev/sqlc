/* name: SelectUserByID :many */
SELECT first_name from
users where (sqlc.arg(id) = id OR sqlc.arg(id) = 0);

/* name: SelectUserByName :many */
SELECT first_name
FROM users
WHERE first_name = sqlc.arg(name)
   OR last_name = sqlc.arg(name);

/* name: SelectUserQuestion :many */
SELECT first_name from
users where (? = id OR  ? = 0);

/* name: SelectUserByAgeCast :many */
SELECT first_name FROM users
WHERE age > CAST(sqlc.arg(threshold) AS SIGNED)
   OR age < CAST(sqlc.arg(threshold) AS SIGNED);
