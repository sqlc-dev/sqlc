CREATE TABLE bars (id serial not null, name text not null, primary key (id));
CREATE TABLE my_data (id serial not null, name text not null, primary key (id));
CREATE TABLE exclusions (id serial not null, name text not null, primary key (id));

-- name: DeleteBarByID :one
DELETE FROM bars WHERE id = $1 RETURNING id, name;

-- name: DeleteMyDataByID :one
DELETE FROM my_data WHERE id = $1 RETURNING id, name;

-- name: DeleteExclusionByID :one
DELETE FROM exclusions WHERE id = $1 RETURNING id, name;
