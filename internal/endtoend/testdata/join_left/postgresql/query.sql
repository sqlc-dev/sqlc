--- https://github.com/kyleconroy/sqlc/issues/604
CREATE TABLE users (
  user_id    INT PRIMARY KEY,
  city_id    INT -- nullable
);
CREATE TABLE cities (
  city_id    INT PRIMARY KEY,
  mayor_id   INT NOT NULL
);
CREATE TABLE mayors (
  mayor_id   INT PRIMARY KEY,
  full_name  TEXT NOT NULL
);

-- name: GetMayors :many
SELECT
    user_id,
    mayors.full_name
FROM users
LEFT JOIN cities USING (city_id)
INNER JOIN mayors USING (mayor_id);

-- name: GetMayorsOptional :many
SELECT
    user_id,
    mayors.full_name
FROM users
LEFT JOIN cities USING (city_id)
LEFT JOIN mayors USING (mayor_id);

-- https://github.com/kyleconroy/sqlc/issues/1334
CREATE TABLE authors (
  id        INT PRIMARY KEY,
  name      TEXT NOT NULL,
  parent_id INT -- nullable
);

CREATE TABLE super_authors (
  super_id        INT PRIMARY KEY,
  super_name      TEXT NOT NULL,
  super_parent_id INT -- nullable
);

-- name: AllAuthors :many
SELECT  *
FROM    authors a
        LEFT JOIN authors p
            ON a.parent_id = p.id;

-- name: AllAuthorsAliases :many
SELECT  *
FROM    authors a
        LEFT JOIN authors p
            ON a.parent_id = p.id;

-- name: AllAuthorsAliases2 :many
SELECT  a.*, p.*
FROM    authors a
        LEFT JOIN authors p
            ON a.parent_id = p.id;

-- name: AllSuperAuthors :many
SELECT  *
FROM    authors
        LEFT JOIN super_authors
            ON authors.parent_id = super_authors.super_id;

-- name: AllSuperAuthorsAliases :many
SELECT  *
FROM    authors a
        LEFT JOIN super_authors sa
            ON a.parent_id = sa.super_id;

-- name: AllSuperAuthorsAliases2 :many
SELECT  a.*, sa.*
FROM    authors a
        LEFT JOIN super_authors sa
            ON a.parent_id = sa.super_id;
