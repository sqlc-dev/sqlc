CREATE TABLE "user" (id bigserial not null);

-- name: User :many
SELECT "user".* FROM "user";
