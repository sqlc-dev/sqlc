CREATE TYPE point_input AS (
    x integer,
    y integer
);

CREATE FUNCTION nearest_to(p point_input)
    RETURNS TABLE
            (
                x integer,
                y integer
            )
    LANGUAGE sql
    STABLE
AS
$$
SELECT p.x, p.y
$$;
