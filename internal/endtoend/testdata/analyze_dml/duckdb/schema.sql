CREATE TABLE authors (
    id BIGINT PRIMARY KEY,
    name TEXT NOT NULL,
    bio TEXT
);

CREATE TABLE books (
    id BIGINT PRIMARY KEY,
    author_id BIGINT NOT NULL,
    title TEXT NOT NULL,
    price DECIMAL(10,2)
);
