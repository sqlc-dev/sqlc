CREATE TABLE users (
  id   INT64  NOT NULL,
  name STRING NOT NULL,
  bio  STRING,
) PRIMARY KEY (id);

CREATE TABLE posts (
  id      INT64     NOT NULL,
  user_id INT64     NOT NULL,
  title   STRING(255),
  created TIMESTAMP NOT NULL,
) PRIMARY KEY (id);
