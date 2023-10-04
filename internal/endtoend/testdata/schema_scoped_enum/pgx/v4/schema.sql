CREATE SCHEMA foo;

CREATE TYPE foo.type_user_role AS ENUM ('admin', 'user');

CREATE TABLE foo.users (
    role foo.type_user_role
);

