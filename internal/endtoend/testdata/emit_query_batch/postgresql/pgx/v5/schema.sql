CREATE SCHEMA myschema;
CREATE TABLE myschema.users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    created_at timestamptz
);
