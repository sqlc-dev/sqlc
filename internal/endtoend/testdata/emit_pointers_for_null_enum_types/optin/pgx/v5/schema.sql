CREATE TYPE user_role AS ENUM ('admin', 'user');

CREATE TABLE users (
    role user_role,
    required_role user_role NOT NULL,
    name text
);
