/* Multi-line
   comment with ; */
CREATE FUNCTION test() RETURNS text AS $$
BEGIN
    -- Internal comment
    RETURN 'test;value';
END;
$$ LANGUAGE plpgsql;