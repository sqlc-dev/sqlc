CREATE FUNCTION plus(a integer, b integer) RETURNS integer AS $$
    BEGIN
        RETURN a + b;
    END;
$$ LANGUAGE plpgsql;

-- name: Plus :one
SELECT plus(b => $2, a => $1);

-- name: MakeIntervalSecs :one
SELECT make_interval(secs => $1);

-- name: MakeIntervalDays :one
SELECT make_interval(days => $1::int);

-- name: MakeIntervalMonths :one
SELECT make_interval(months => sqlc.arg('months')::int);


