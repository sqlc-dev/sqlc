/* name: TooManyFroms :one */
select id, first_name from users from where id = $1;
