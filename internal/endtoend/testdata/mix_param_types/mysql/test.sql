CREATE TABLE bar (
       id serial not null,
       name text not null,
       phone text not null
);

-- name: CountOne :one
SELECT count(1) FROM bar WHERE id = sqlc.arg(id) AND name <> ?; 

-- name: CountTwo :one
SELECT count(1) FROM bar WHERE id = ? AND name <> sqlc.arg(name);

-- name: CountThree :one
SELECT count(1) FROM bar WHERE id > ? AND phone <> sqlc.arg(phone) AND name <> ?;
