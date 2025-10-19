CREATE FUNCTION foo() RETURNS text AS $$
BEGIN
    RETURN 'test;';
END;
$$ LANGUAGE plpgsql;