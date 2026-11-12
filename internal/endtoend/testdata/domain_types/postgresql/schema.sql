CREATE DOMAIN tags AS text[];
CREATE DOMAIN title AS varchar(100) NOT NULL DEFAULT '';
CREATE DOMAIN price AS numeric(10, 2) CHECK (VALUE >= 0);
CREATE SCHEMA billing;
CREATE DOMAIN billing.ksuid AS char(27);
-- Named like a built-in, but off the search path.
CREATE DOMAIN billing.text AS integer;
CREATE DOMAIN dropped AS text;
DROP DOMAIN dropped;

CREATE TABLE books (
  id serial PRIMARY KEY,
  title title,
  price price,
  ref billing.ksuid
);

CREATE TABLE shelves (
  id serial PRIMARY KEY,
  titles title[] NOT NULL,
  tags tags NOT NULL,
  tag_sets tags[] NOT NULL,
  label title,
  note text
);
