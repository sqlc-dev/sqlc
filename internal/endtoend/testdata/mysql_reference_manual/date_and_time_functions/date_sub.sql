-- name: DateSubOneYear :one
SELECT DATE_SUB('2018-05-01',INTERVAL 1 YEAR);

-- name: DateSubDaySecond :one
SELECT DATE_SUB('2025-01-01 00:00:00',
                INTERVAL '1 1:1:1' DAY_SECOND);

-- name: DateSub31Days :one
SELECT DATE_SUB('1998-01-02', INTERVAL 31 DAY);
