-- name: User :many
SELECT "user".* FROM "user";

-- name: UsersB :batchmany
SELECT * FROM "user"
WHERE id = $1;

-- name: UsersC :copyfrom
INSERT INTO "user"
(id)
VALUES
    ($1);
