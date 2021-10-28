CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- name: WordSimilarity :one
SELECT word_similarity('word', 'two words');
