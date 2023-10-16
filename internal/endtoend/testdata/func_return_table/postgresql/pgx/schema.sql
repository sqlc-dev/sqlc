CREATE TABLE accounts (
  id         INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  username   TEXT NOT NULL UNIQUE,
  password   TEXT NOT NULL
);

-- this is a useless and horrifying function cause we don't hash
-- the password, this is just to repro the bug in sqlc
CREATE OR REPLACE FUNCTION register_account(
    _username TEXT,
    _password VARCHAR(70)
)
RETURNS TABLE (
    account_id   INTEGER
)
AS $$
BEGIN
  INSERT INTO accounts (username, password)
       VALUES (
         _username,
         _password
       )
    RETURNING id INTO account_id;

  RETURN NEXT;
END;
$$ LANGUAGE plpgsql;
