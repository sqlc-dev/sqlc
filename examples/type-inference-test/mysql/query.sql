-- name: ListTest :many
SELECT 
    (a1 / 1024) a1_float, (a2 / 1024) a2_float, a3
FROM test;

-- name: ListTest2 :many
SELECT 
    COALESCE(CAST(a1 / 1024 AS FLOAT), 0) a1_float, COALESCE(CAST(a2 / 1024 AS FLOAT), 0) a2_float, a3
FROM test;

-- name: ListTest3 :many
SELECT 
    CAST(a1 / 1024 AS FLOAT) a1_float, CAST(a2 / 1024 AS FLOAT) a2_float, a3
FROM test;

-- name: ListTest4 :many
SELECT 
    (a1 + a2) as sum_result,
    (a1 * a2) as mult_result,
    (a1 - a2) as sub_result,
    (a1 % 10) as mod_result
FROM test;

-- name: ListTest5 :many
SELECT 
    COALESCE(a1 / 1024, 0) as with_inference,
    COALESCE(a2 / 1024, 0) as nullable_inference
FROM test;
