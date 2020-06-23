CREATE TYPE job_status AS ENUM ('APPLIED', 'PENDING', 'ACCEPTED', 'REJECTED');

CREATE TABLE users (
    id integer NOT NULL,
    first_name varchar(255) NOT NULL,
    last_name varchar(255),
    age integer NOT NULL,
    job_status job_status NOT NULL
);

CREATE TABLE orders (
    id integer NOT NULL,
    price DECIMAL(13, 4) NOT NULL,
    user_id integer NOT NULL
);


