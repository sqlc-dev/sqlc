CREATE TABLE authors (
  id         bigserial PRIMARY KEY,
  org_id     bigint NOT NULL,
  name       text NOT NULL,
  deleted_at timestamptz
);
