CREATE TABLE authors (
    author_id SERIAL PRIMARY KEY
);

-- name: CreateAuthors :copyfrom
INSERT INTO authors (author_id) VALUES ($1);
