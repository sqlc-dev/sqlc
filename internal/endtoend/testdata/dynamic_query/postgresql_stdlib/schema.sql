CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    status INT -- 0: inactive, 1: active, 2: pending
);
