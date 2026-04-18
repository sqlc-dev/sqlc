CREATE TABLE users (
    id integer NOT NULL AUTO_INCREMENT PRIMARY KEY,
    first_name varchar(255) NOT NULL,
    last_name varchar(255),
    age integer NOT NULL,
    shoe_size ENUM('x-small', 'small', 'medium', 'large', 'x-large') NOT NULL,
    shirt_size ENUM('x-small', 'small', 'medium', 'large', 'x-large')
);
