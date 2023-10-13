-- name: CreateAuthors :copyfrom
INSERT INTO authors (author_id) VALUES ($1);
