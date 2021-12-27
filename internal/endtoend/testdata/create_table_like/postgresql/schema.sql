CREATE TABLE users (
    id serial NOT NULL PRIMARY KEY,
    first_name varchar(255) NOT NULL
);

ALTER TABLE users
    ADD COLUMN last_name varchar(255);

CREATE TABLE super_users ( LIKE users );

ALTER TABLE users
    ADD COLUMN age integer NOT NULL;

ALTER TABLE users
    DROP COLUMN first_name;

ALTER TABLE super_users
    ADD COLUMN date_of_birth TIMESTAMP;
