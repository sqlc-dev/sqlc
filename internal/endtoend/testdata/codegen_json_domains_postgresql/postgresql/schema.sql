CREATE DOMAIN title AS varchar(100) NOT NULL DEFAULT '';
CREATE DOMAIN price AS numeric(10, 2) CHECK (VALUE >= 0);
CREATE SCHEMA billing;
CREATE DOMAIN billing.ksuid AS char(27);
CREATE DOMAIN dropped AS text;
DROP DOMAIN dropped;

CREATE TABLE books (
  id serial PRIMARY KEY,
  title title,
  price price,
  ref billing.ksuid
);
COMMENT ON DOMAIN title IS 'A book''s title';
