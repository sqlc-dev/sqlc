CREATE TABLE users (
    id SERIAL NOT NULL,
    first_name varchar(255) NOT NULL,
    tags text[] NOT NULL
);
