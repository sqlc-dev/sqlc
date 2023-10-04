CREATE TABLE users (
    id integer NOT NULL AUTO_INCREMENT PRIMARY KEY,
    first_name varchar(255) NOT NULL,
    last_name varchar(255),
    age integer NOT NULL,
    job_status varchar(10) NOT NULL
) ENGINE=InnoDB;

CREATE TABLE orders (
    id integer NOT NULL AUTO_INCREMENT PRIMARY KEY,
    price DECIMAL(13, 4) NOT NULL,
    user_id integer NOT NULL
) ENGINE=InnoDB;

