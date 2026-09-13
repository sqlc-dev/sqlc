CREATE TABLE users (
  id   INT64  NOT NULL,
  name STRING(MAX) NOT NULL,
  bio  STRING(MAX),
) PRIMARY KEY (id);

CREATE TABLE posts (
  id      INT64     NOT NULL,
  user_id INT64     NOT NULL,
  title   STRING(255),
  created TIMESTAMP NOT NULL,
) PRIMARY KEY (id);
