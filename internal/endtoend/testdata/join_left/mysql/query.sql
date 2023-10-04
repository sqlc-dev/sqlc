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
    cities.city_id,
    mayors.full_name
FROM users
LEFT JOIN cities USING (city_id)
LEFT JOIN mayors USING (mayor_id);

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

-- name: GetSuggestedUsersByID :many
SELECT  DISTINCT u.*, m.*
FROM    users_2 u
        LEFT JOIN media m
            ON u.user_avatar_id = m.media_id
WHERE   u.user_id != @user_id;
