CREATE SCHEMA old;
CREATE SCHEMA new;

CREATE TYPE event AS enum ('START', 'STOP');
CREATE TYPE old.level AS enum ('DEBUG', 'INFO', 'WARN', 'ERROR', 'FATAL');

CREATE TABLE log_lines (
  id     BIGSERIAL    PRIMARY KEY,
  status event        NOT NULL,
  level  old.level    NOT NULL
);

ALTER TYPE event SET SCHEMA new;
ALTER TYPE old.level SET SCHEMA public;

