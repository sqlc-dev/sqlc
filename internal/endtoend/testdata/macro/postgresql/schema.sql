CREATE TABLE authors (
  id         bigserial PRIMARY KEY,
  org_id     bigint NOT NULL,
  name       text NOT NULL,
  bio        text,
  created_at timestamptz NOT NULL,
  deleted_at timestamptz
);
