CREATE DOMAIN tags AS text[];
CREATE DOMAIN title AS varchar(100) NOT NULL DEFAULT '';
CREATE DOMAIN price AS numeric(10, 2) CHECK (VALUE >= 0);
-- The base keeps a quoted name's case and pg_catalog.
CREATE TYPE "Mood" AS ENUM ('happy', 'sad');
CREATE DOMAIN feeling AS "Mood";
CREATE DOMAIN amount AS numeric(10, 2) NOT NULL;
CREATE SCHEMA billing;
CREATE DOMAIN billing.ksuid AS char(27);
-- Named like a built-in, but off the search path.
CREATE DOMAIN billing.text AS integer;
CREATE DOMAIN relaxed AS text NOT NULL;
CREATE DOMAIN tightened AS text;
-- A domain carries the NOT NULL of the domain it stands on, an array
-- of that domain does not.
CREATE DOMAIN required_text AS text NOT NULL;
CREATE DOMAIN derived AS required_text;
CREATE DOMAIN derived_list AS required_text[];
CREATE DOMAIN dropped AS text;
DROP DOMAIN dropped;

CREATE TABLE books (
  id serial PRIMARY KEY,
  title title,
  price price,
  ref billing.ksuid,
  feeling feeling NOT NULL,
  amount amount
);

CREATE TABLE shelves (
  id serial PRIMARY KEY,
  titles title[] NOT NULL,
  tags tags NOT NULL,
  tag_sets tags[] NOT NULL,
  label title,
  note text,
  retyped title,
  strict title NOT NULL,
  relaxed relaxed,
  tightened tightened,
  derived derived,
  derived_list derived_list
);

-- A child's own NOT NULL outlasts the domain's.
CREATE TABLE notes (body relaxed);
CREATE TABLE pinned_notes (body relaxed NOT NULL) INHERITS (notes);

-- The domain's NOT NULL leaves with it, and outlasts the column's own.
ALTER TABLE shelves ALTER COLUMN retyped TYPE text;
ALTER TABLE shelves ALTER COLUMN strict DROP NOT NULL;

ALTER DOMAIN relaxed DROP NOT NULL;
ALTER DOMAIN tightened SET NOT NULL;
