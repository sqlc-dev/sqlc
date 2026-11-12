CREATE DOMAIN tags AS text[];
CREATE DOMAIN title AS varchar(100) NOT NULL;

CREATE TABLE shelves (
  id serial PRIMARY KEY,
  titles title[] NOT NULL,
  tags tags NOT NULL,
  tag_sets tags[] NOT NULL
);
