-- Example queries for sqlc
CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL,
  bio  text,
  country_code CHAR(2) NOT NULL,
  titles TEXT[]
);

CREATE TABLE clients (
  id INT PRIMARY KEY,
  name TEXT NOT NULL
);
