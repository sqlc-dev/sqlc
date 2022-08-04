-- name: GetAvg :one
SELECT avg(int_val) FROM test;

-- name: GetCount :one
SELECT count(*) FROM test;

-- name: GetCountId :one
SELECT count(id) FROM test;

-- name: GetGroupConcatInt :one
SELECT group_concat(int_val) FROM test;

-- name: GetGroupConcatInt2 :one
SELECT group_concat(1, ':') FROM test;

-- name: GetGroupConcatText :one
SELECT group_concat(text_val) FROM test;

-- name: GetGroupConcatText2 :one
SELECT group_concat(text_val, ':') FROM test;

-- name: GetMaxInt :one
SELECT max(int_val) FROM test;

-- name: GetMaxText :one
SELECT max(text_val) FROM test;

-- name: GetMinInt :one
SELECT min(int_val) FROM test;

-- name: GetMinText :one
SELECT min(text_val) FROM test;

-- name: GetSumInt :one
SELECT sum(int_val) FROM test;

-- name: GetSumText :one
SELECT sum(text_val) FROM test;

-- name: GetTotalInt :one
SELECT total(int_val) FROM test;
