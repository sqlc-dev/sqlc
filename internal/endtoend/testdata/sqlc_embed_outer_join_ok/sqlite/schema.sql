CREATE TABLE users (
    id integer PRIMARY KEY,
    name text NOT NULL
);

CREATE TABLE posts (
    id integer PRIMARY KEY,
    user_id integer
);
