CREATE TYPE user_role AS ENUM ('admin', 'user');

CREATE SCHEMA foo;

CREATE TYPE foo.status AS ENUM ('active', 'inactive');

CREATE TABLE users (
    role user_role,
    required_role user_role NOT NULL,
    status foo.status
);
