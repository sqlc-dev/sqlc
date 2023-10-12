-- name: StageUserData :copyfrom
insert into "user_data" ("id", "user")
values (
    sqlc.arg('id_param'),
    sqlc.arg('user_param')
);
