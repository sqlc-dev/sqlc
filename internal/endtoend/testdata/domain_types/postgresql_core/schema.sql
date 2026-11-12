CREATE DOMAIN title AS varchar(100) NOT NULL;

CREATE TABLE shelves (
  id serial PRIMARY KEY,
  titles title[] NOT NULL
);
