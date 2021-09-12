-- https://dev.mysql.com/doc/refman/8.0/en/date-and-time-functions.html#function_date-add

-- name: DateAddOneDay :one
SELECT DATE_ADD('2018-05-01',INTERVAL 1 DAY);

-- name: DateAddOneSecond :one
SELECT DATE_ADD('2020-12-31 23:59:59',
                INTERVAL 1 SECOND);

-- name: DateAddTimestampOneSecond :one
SELECT DATE_ADD('2018-12-31 23:59:59',
                INTERVAL 1 DAY);

-- name: DateAddMinuteSecond :one
SELECT DATE_ADD('2100-12-31 23:59:59',
                INTERVAL '1:1' MINUTE_SECOND);

-- name: DateAddDayHour :one
SELECT DATE_ADD('1900-01-01 00:00:00',
                INTERVAL '-1 10' DAY_HOUR);

-- name: DateAddSecondMicrosecond :one
SELECT DATE_ADD('1992-12-31 23:59:59.000002',
           INTERVAL '1.999999' SECOND_MICROSECOND);
