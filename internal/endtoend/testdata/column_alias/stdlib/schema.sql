CREATE EXTENSION pg_trgm;
CREATE EXTENSION pgcrypto;

CREATE TABLE users(
    id INT GENERATED ALWAYS AS IDENTITY NOT NULL,
    fname VARCHAR(100) NOT NULL,
    lname VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    enc_passwd TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL default (NOW() AT TIME ZONE 'utc'),
    PRIMARY KEY(id)
);
