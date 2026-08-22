CREATE TABLE authors (
  id BIGSERIAL PRIMARY KEY,
  name text NOT NULL,
  bio text,
  created_at timestamptz NOT NULL DEFAULT now()
);
