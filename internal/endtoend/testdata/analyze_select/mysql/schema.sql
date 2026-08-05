CREATE TABLE users (
  id   BIGINT       PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  bio  TEXT
);

CREATE TABLE posts (
  id      BIGINT   PRIMARY KEY,
  user_id BIGINT   NOT NULL,
  title   VARCHAR(255),
  created DATETIME NOT NULL
);
