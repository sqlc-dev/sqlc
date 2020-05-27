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
