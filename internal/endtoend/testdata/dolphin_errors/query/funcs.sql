/* name: WrongFunc :one */
select id, first_name from users where id = SQLC_ARGH(target_id);

/* name: InvalidName :one */
select id, first_name from users where id = SQLC_ARG(SQLC_ARG(target_id));

/* name: InvalidVaue :one */
select id, first_name from users where id = SQLC_ARG(?);
