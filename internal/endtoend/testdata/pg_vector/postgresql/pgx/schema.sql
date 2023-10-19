CREATE EXTENSION IF NOT EXISTS "vector";

CREATE TABLE items (id bigserial PRIMARY KEY, embedding vector(3));
