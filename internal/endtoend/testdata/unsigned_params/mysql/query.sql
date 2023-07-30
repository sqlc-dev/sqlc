CREATE TABLE foo (id INT UNSIGNED NOT NULL);

-- name: CreateFoo :exec
INSERT INTO foo (id) VALUES (?);
