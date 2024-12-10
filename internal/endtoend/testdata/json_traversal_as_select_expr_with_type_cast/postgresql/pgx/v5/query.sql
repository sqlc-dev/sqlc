-- name: MyGet :many
SELECT *,(mt.myjson -> 'thing1' -> 'thing2')::text,mt.myjson -> 'thing1'
FROM mytable mt;
