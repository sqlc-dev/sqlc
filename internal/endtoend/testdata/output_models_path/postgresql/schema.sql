CREATE TYPE status AS ENUM ('active', 'inactive', 'pending');

CREATE TABLE authors (
  id     BIGSERIAL PRIMARY KEY,
  name   text   NOT NULL,
  bio    text,
  status status NOT NULL DEFAULT 'active'
);

CREATE TABLE books (
  id        BIGSERIAL PRIMARY KEY,
  author_id BIGINT NOT NULL REFERENCES authors(id),
  title     text   NOT NULL
);
