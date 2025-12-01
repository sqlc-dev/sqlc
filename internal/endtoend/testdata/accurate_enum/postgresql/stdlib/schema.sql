CREATE TYPE status AS ENUM ('pending', 'active', 'completed');

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    status status NOT NULL DEFAULT 'pending'
);
