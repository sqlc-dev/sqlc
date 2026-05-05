CREATE TYPE status AS ENUM ('active', 'inactive');

CREATE TABLE authors (
  id     BIGSERIAL PRIMARY KEY,
  name   text   NOT NULL,
  status status NOT NULL DEFAULT 'active'
);
