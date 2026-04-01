CREATE TABLE servers (
  code    varchar PRIMARY KEY,
  name    text    NOT NULL,
  count   integer NOT NULL DEFAULT 0
);
