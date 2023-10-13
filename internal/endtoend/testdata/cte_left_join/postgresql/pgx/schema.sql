CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL,
  bio  text
);

CREATE TABLE fake (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL,
  bio  text
);