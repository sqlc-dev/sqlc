CREATE TYPE event AS enum ('START', 'STOP');

CREATE TABLE log_lines (
  id     BIGSERIAL    PRIMARY KEY,
  status "event"  NOT NULL
);

ALTER TYPE event RENAME TO "new_event";

-- name: ListAuthors :many
SELECT * FROM log_lines;
