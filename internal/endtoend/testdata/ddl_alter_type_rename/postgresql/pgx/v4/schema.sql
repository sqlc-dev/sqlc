CREATE TYPE event AS enum ('START', 'STOP');

ALTER TYPE event RENAME TO "new_event";

CREATE TABLE log_lines (
  id     BIGSERIAL    PRIMARY KEY,
  status "new_event"  NOT NULL
);

