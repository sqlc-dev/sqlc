CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- name: GenerateUUID :one
SELECT uuid_generate_v4();
