CREATE TABLE authors (
    id BIGINT PRIMARY KEY,
    name TEXT NOT NULL,
    bio TEXT,
    royalties DECIMAL(10,2) NOT NULL
);
