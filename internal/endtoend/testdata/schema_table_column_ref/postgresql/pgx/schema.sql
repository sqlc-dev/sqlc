CREATE SCHEMA astoria;

CREATE TABLE astoria.slack_feedback (
  id            BIGSERIAL PRIMARY KEY,
  workspace_id  BIGINT NOT NULL, 
  created_at    TIMESTAMP NOT NULL,
  issue_raised  BOOLEAN
);

CREATE TABLE astoria.tickets (
  id            BIGSERIAL PRIMARY KEY,
  workspace_id  BIGINT NOT NULL, 
  created_at    TIMESTAMP NOT NULL,
  source text   NOT NULL
);
