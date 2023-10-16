-- Example queries for sqlc
CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL,
  bio  text
);

CREATE OR REPLACE FUNCTION add_author (name text, bio text, out id int)
AS $$
DECLARE
BEGIN
  id = 123;
END;
$$ LANGUAGE plpgsql;
