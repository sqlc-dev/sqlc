CREATE TABLE bar (
	id serial not null,
	name text not null,
	phone text not null
);

-- name: CountOne :one
SELECT count(1) FROM bar WHERE id = $2 AND phone < @phone_param and name <> $1;

-- name: CountTwo :one
SELECT count(1) FROM bar WHERE id = sqlc.arg(id_param) AND phone < @phone_param and name <> $1;

-- name: CountThree :one
SELECT count(1) FROM bar WHERE id > sqlc.arg(id_param) AND name = $1;

-- name: CountFour :one
SELECT count(1) FROM bar WHERE id > $2 AND phone <> @phone_param  AND name <> $1;
