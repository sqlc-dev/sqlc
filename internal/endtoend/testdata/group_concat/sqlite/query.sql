-- name: GroupConcatOne :one
SELECT group_concat(name order by id asc) FROM book;

-- name: GroupConcat :one
SELECT group_concat(name order by id asc) FROM book;

-- name: GroupConcatDelimeter :one
SELECT group_concat(name, ',' order by id asc) FROM book;

-- name: StringAgg :one
SELECT string_agg(name, ',' order by id asc) FROM book;
