CREATE FUNCTION plus(a integer, b integer) RETURNS integer AS $$
    BEGIN
        RETURN a + b;
    END;
$$ LANGUAGE plpgsql;

CREATE FUNCTION table_args(x INT) RETURNS TABLE (y INT) AS 'SELECT x' LANGUAGE sql;

