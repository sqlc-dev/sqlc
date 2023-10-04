-- https://github.com/sqlc-dev/sqlc/issues/1728

CREATE TABLE IF NOT EXISTS locations (
    id              INTEGER PRIMARY KEY,
    name            TEXT    NOT NULL,
    address         TEXT    NOT NULL,
    zip_code        INT     NOT NULL,
    latitude        REAL    NOT NULL,
    longitude       REAL    NOT NULL,
    UNIQUE(name)
);

