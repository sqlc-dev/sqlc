CREATE TYPE pair AS (a integer, b integer);

CREATE FUNCTION sum_pairs(inputs pair[])
    RETURNS TABLE (total bigint)
    LANGUAGE sql
    STABLE
AS
$$
SELECT 0::bigint
$$;
