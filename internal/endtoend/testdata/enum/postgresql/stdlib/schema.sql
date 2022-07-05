CREATE TYPE size AS ENUM('x-small', 'small', 'medium', 'large', 'x-large');

CREATE TABLE users (
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    first_name varchar(255) NOT NULL,
    last_name varchar(255),
    age integer NOT NULL,
    shoe_size size NOT NULL,
    shirt_size size
);
