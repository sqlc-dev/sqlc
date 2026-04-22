-- name: GetUsers :many
select * from "user" where is_deleted = false;
