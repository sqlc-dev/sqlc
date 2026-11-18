-- https://github.com/sqlc-dev/sqlc/issues/3440
CREATE TABLE t (
    id      INTEGER PRIMARY KEY,
    name    TEXT    NOT NULL,
    score   INTEGER NOT NULL,
    created TEXT
);

CREATE UNIQUE INDEX t_name ON t (name);
