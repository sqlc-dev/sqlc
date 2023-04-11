CREATE TABLE "user" (id bigserial not null);

-- name: User :many
SELECT "user".* FROM "user";

-- name: UsersB :batchmany
SELECT * FROM "user"
WHERE id = $1;
