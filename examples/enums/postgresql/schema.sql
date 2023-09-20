CREATE TYPE size AS ENUM('x-small', 'small', 'medium', 'large', 'x-large');

CREATE TABLE users (
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    first_name text NOT NULL,
    last_name text NOT NULL,
    age integer NOT NULL,
    shirt_size size NOT NULL
);
