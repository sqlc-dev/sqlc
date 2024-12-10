-- name: MyGet :many
SELECT *, (mt.myjson->'thing1'->'thing2')::text
FROM mytable mt;
