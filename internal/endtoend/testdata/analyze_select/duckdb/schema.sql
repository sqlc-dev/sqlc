CREATE TABLE users (
  id   BIGINT NOT NULL,
  name TEXT   NOT NULL,
  age  INTEGER
);

CREATE TABLE posts (
  id      BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  title   TEXT   NOT NULL,
  body    TEXT
);
