CREATE TABLE authors (
  id         bigserial PRIMARY KEY,
  name       text NOT NULL,
  deleted_at timestamptz
);

CREATE TABLE books (
  id         bigserial PRIMARY KEY,
  author_id  bigint NOT NULL REFERENCES authors (id),
  title      text NOT NULL,
  deleted_at timestamptz
);
