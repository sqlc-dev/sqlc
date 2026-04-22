-- name: GetUsers :many
select * from "user" where exists (
  select 1 from "user" as u2 where is_deleted = false
);
