CREATE TABLE users (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL,
  bio  text
);

CREATE TABLE posts (
  id      BIGSERIAL PRIMARY KEY,
  user_id bigint      NOT NULL,
  title   varchar(255),
  created timestamptz NOT NULL
);
