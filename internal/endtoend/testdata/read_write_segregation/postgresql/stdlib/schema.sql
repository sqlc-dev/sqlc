CREATE TABLE users (
    id SERIAL NOT NULL,
    first_name varchar(255) NOT NULL,
    last_name varchar(255),
    updated_at timestamp NULL DEFAULT CURRENT_TIMESTAMP
);

