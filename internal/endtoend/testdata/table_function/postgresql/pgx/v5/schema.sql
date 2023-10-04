CREATE TABLE transactions (
  id  BIGSERIAL PRIMARY KEY,
  uri text NOT NULL,
  program_id text NOT NULL,
  data text NOT NULL
);

