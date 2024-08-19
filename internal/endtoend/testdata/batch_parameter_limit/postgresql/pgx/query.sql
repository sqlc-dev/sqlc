-- name: CreateAuthors :batchexec
INSERT INTO authors (
  name, bio
) VALUES (
  $1, $2
);
