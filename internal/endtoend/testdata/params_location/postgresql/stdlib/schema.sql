CREATE TABLE users (
    id SERIAL NOT NULL,
    first_name varchar(255) NOT NULL,
    last_name varchar(255),
    age integer NOT NULL,
    job_status varchar(10) NOT NULL
);

CREATE TABLE orders (
    id SERIAL NOT NULL,
    price DECIMAL(13, 4) NOT NULL,
    user_id integer NOT NULL
);

