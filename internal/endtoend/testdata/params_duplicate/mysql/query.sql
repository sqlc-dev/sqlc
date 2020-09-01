/* name: SelectUserArg :many */
SELECT first_name from
users where (sqlc.arg(id) = id OR sqlc.arg(id) = 0);

/* name: SelectUserQuestion :many */
SELECT first_name from
users where (? = id OR  ? = 0);
