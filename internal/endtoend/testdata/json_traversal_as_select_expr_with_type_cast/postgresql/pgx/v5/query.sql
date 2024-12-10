-- name: MyGet :many
SELECT *, (mt.myjson->'thing1'->'thing2')::text, mt.myjson->'thing1'
FROM "mytable" mt;

-- name: MyGet2 :many
SELECT id::text
FROM "mytable" mt;