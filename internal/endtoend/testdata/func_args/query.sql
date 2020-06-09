CREATE FUNCTION plus(a integer, b integer) RETURNS integer AS $$
    BEGIN
        RETURN a + b;
    END;
$$ LANGUAGE plpgsql;

CREATE FUNCTION table_args(x INT) RETURNS TABLE (y INT) AS 'SELECT x' LANGUAGE sql;

-- name: Plus :one
SELECT plus(b => $2, a => $1);

-- name: MakeIntervalSecs :one
SELECT make_interval(secs => $1);

-- name: MakeIntervalDays :one
SELECT make_interval(days => $1::int);

-- name: MakeIntervalMonths :one
SELECT make_interval(months => sqlc.arg('months')::int);

-- name: TableArgs :one
SELECT table_args(x => $1);
