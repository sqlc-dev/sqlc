CREATE DOMAIN tags AS text[];
CREATE DOMAIN title AS varchar(100) NOT NULL;
-- An override on the domain a domain stands on applies to it.
CREATE DOMAIN id AS text;
CREATE DOMAIN validated_id AS id;

CREATE TABLE shelves (
  id serial PRIMARY KEY,
  titles title[] NOT NULL,
  tags tags NOT NULL,
  tag_sets tags[] NOT NULL
);

CREATE TABLE things (
  ref validated_id NOT NULL,
  refs validated_id[] NOT NULL
);
