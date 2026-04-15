-- name: InsertUser :exec
INSERT INTO users (name) VALUES (?);

-- name: InsertUserMixedCase :exec
INSERT INTO users (name) VALUES (?);

-- name: InsertAuthor :exec
INSERT INTO "Authors" (name) VALUES (?);

-- name: InsertBook :exec
INSERT INTO Books (title) VALUES (?);

-- name: UpdateUser :exec
UPDATE users SET name = ? WHERE id = ?;

-- name: UpdateUserMixedCase :exec
UPDATE users SET name = ? WHERE id = ?;

-- name: UpdateAuthor :exec
UPDATE "Authors" SET name = ? WHERE id = ?;

-- name: UpdateBook :exec
UPDATE Books SET title = ? WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;

-- name: DeleteUserMixedCase :exec
DELETE FROM users WHERE id = ?;

-- name: DeleteAuthor :exec
DELETE FROM "Authors" WHERE id = ?;

-- name: DeleteBook :exec
DELETE FROM Books WHERE id = ?;

-- name: GetUser :one
SELECT * FROM users WHERE id = ?;

-- name: GetUserMixedCase :one
SELECT * FROM users WHERE id = ?;

-- name: GetAuthor :one
SELECT * FROM "Authors" WHERE id = ?;

-- name: GetBook :one
SELECT * FROM Books WHERE id = ?;
