-- Example queries for sqlc
CREATE TABLE authors (
    id BigSerial,
    name Text NOT NULL,
    bio Text,
    PRIMARY KEY (id)
);

-- https://github.com/sqlc-dev/sqlc/issues/1290
CREATE TABLE users (
    sub Uuid,
    PRIMARY KEY (sub)
);
