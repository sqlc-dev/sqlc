-- name: GetFullName :many
SELECT id, concat(first_name, ' ', last_name) AS full_name FROM employees;

-- name: GetUppercaseNames :many
SELECT id, upper(first_name) AS first_name_upper, lower(last_name) AS last_name_lower FROM employees;

-- name: GetEmailDomain :many
SELECT id, email, substring(email, position(email, '@') + 1) AS domain FROM employees;

-- name: TrimWhitespace :many
SELECT id, trim(bio) AS bio FROM employees;
