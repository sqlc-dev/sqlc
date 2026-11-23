CREATE TABLE authors (
  id         integer PRIMARY KEY,
  org_id     integer NOT NULL,
  name       text NOT NULL,
  bio        text,
  deleted_at datetime
);
