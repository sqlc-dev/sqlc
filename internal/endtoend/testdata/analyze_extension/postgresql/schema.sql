CREATE EXTENSION citext;
CREATE EXTENSION pg_trgm;

CREATE TABLE users (
  id    BIGSERIAL PRIMARY KEY,
  name  text   NOT NULL,
  email citext NOT NULL
);
