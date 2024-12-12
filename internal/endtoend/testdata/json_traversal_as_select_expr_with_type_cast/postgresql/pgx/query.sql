-- name: GetNullable1 :many
SELECT null::text
FROM "mytable";

-- name: GetNullable2 :many
SELECT CASE
    WHEN id = 1 THEN id::int
    WHEN id = 2 THEN null
    WHEN id = 3 THEN 8.5
    WHEN id = 4 THEN 7
    ELSE '2'
END
FROM "mytable";

-- name: GetNullable3 :many
SELECT CASE WHEN true THEN 'hello'::text ELSE null END
FROM "mytable";

-- name: GetNullable4 :many
SELECT CASE WHEN true THEN 'hello'::text END
FROM "mytable";

-- name: GetNullable5 :many
SELECT (mt.myjson->'thing1'->'thing2')::text
FROM "mytable" mt;

-- name: GetNullable6 :many
SELECT mt.myjson->'thing1'->>'thing2'
FROM "mytable" mt;

-- name: GetNullable7 :many
SELECT mt.myjson->'thing1'->'thing2'
FROM "mytable" mt;

-- name: GetNullable2A :many
SELECT CASE
    WHEN id = 1 THEN id::int
    WHEN id = 2 THEN null
    WHEN id = 3 THEN 8.5
    ELSE 7
END
FROM "mytable";

-- name: GetNullable2B :many
SELECT CASE
    WHEN id = 1 THEN id::float
    WHEN id = 2 THEN null
    ELSE 7
END
FROM "mytable";

-- name: GetNullable2C :many
SELECT CASE
    WHEN id = 1 THEN true
    ELSE null
    END
FROM "mytable";

-- name: GetNullable2D :many
SELECT CASE
    WHEN id = 2 THEN mt.myjson->'thing1'->>'thing2'
    ELSE null
    END
FROM "mytable" mt;

-- name: GetNullable2E :many
SELECT CASE
    WHEN id = 2 THEN mt.myjson->'thing1'->'thing2'
    WHEN id = 3 THEN mt.myjson->'thing1'
    ELSE null
    END
FROM "mytable" mt;

-- name: GetNullable2F :many
SELECT CASE
    WHEN id = 2 THEN null
    ELSE 7 - id
END
FROM "mytable";