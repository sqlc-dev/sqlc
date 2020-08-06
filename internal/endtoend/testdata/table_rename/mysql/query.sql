CREATE TABLE authors (
    id   BIGINT PRIMARY KEY,
    name TEXT NOT NULL,
    bio  TEXT
);

RENAME TABLE authors TO writers;

SELECT * FROM writers;
