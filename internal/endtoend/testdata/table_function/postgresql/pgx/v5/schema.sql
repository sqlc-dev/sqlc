CREATE TABLE transactions (
  id         BIGSERIAL PRIMARY KEY,
  uri        TEXT      NOT NULL,
  program_id TEXT      NOT NULL,
  data       JSONB      NOT NULL
);

