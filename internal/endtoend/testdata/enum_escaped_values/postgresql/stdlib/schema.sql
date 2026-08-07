CREATE TYPE injection AS ENUM (
    'safe_value',
    'injected" + "arbitrary_go_code" + "'
);

CREATE TYPE user_role AS ENUM (
    'admin',
    'user\nadmin'
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    role user_role NOT NULL,
    payload injection NOT NULL
);
