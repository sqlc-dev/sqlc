CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL
);

CREATE TABLE books (
  id        BIGSERIAL PRIMARY KEY,
  author_id bigint    NOT NULL REFERENCES authors (id),
  title     text      NOT NULL
);
