-- Create a user
INSERT INTO users (name, email) VALUES ('John''s', 'john@example.com'); -- comment;

/* Multi-line
   comment with ; */
CREATE FUNCTION test() RETURNS text AS $$
BEGIN
    -- Internal comment
    RETURN 'test;value';
END;
$$ LANGUAGE plpgsql;

SELECT "weird;column" FROM users WHERE name = 'test;value';