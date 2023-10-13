CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  name text NOT NULL
);

CREATE TABLE books (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL,
  author_id BIGSERIAL REFERENCES authors(id)
);