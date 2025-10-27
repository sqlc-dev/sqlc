-- name: SelectUserByID :many
SELECT first_name from
users where ($id = id OR $id = 0);

-- name: SelectUserByName :many
SELECT first_name
FROM users
WHERE first_name = $name
   OR last_name = $name;

-- name: SelectUserQuestion :many
SELECT first_name from
users where ($question = id OR $question = 0);
