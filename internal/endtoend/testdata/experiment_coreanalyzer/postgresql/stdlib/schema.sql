CREATE TABLE authors (
    id   bigint PRIMARY KEY,
    name text NOT NULL,
    bio  text,
    tags text[]
);
