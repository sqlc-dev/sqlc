CREATE TABLE versions (
  id   BIGSERIAL PRIMARY KEY,
  name TEXT,
  previous_version_id bigint NOT NULL
);
