-- Example queries for sqlc
CREATE TABLE "test"
(
    "id" TEXT NOT NULL
);

-- name: TestList :many
SELECT * FROM "test";