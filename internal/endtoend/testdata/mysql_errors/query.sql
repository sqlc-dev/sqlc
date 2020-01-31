/* name: WrongFunc :one */
select id, first_name from users where id = sqlc.argh(target_id);

/* name: InvalidName :one */
select id, first_name from users where id = sqlc.arg(sqlc.arg(target_id));

/* name: InvalidVaue :one */
select id, first_name from users where id = sqlc.arg(?);

/* name: TooManyFroms :one */
select id, first_name from users from where id = ?;

/* name: MisspelledSelect :one */
selectt id, first_name from users;

/* name: ExtraSelect :one */
select id from users where select id;

-- stderr
-- # package querytest
-- query.sql:1:1: invalid function call "sqlc.argh", did you mean "sqlc.arg"?
-- query.sql:4:1: invalid custom argument value "sqlc.arg(sqlc.arg(target_id))"
-- query.sql:7:1: invalid custom argument value "sqlc.arg(?)"
-- query.sql:11:39: syntax error at or near "from"
-- query.sql:14:9: syntax error at or near "selectt"
-- query.sql:17:35: syntax error at or near "select"
