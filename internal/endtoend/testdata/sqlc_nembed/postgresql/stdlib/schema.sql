CREATE TABLE users (
    id integer NOT NULL PRIMARY KEY,
    name text NOT NULL
);

CREATE TABLE posts (
    id integer NOT NULL PRIMARY KEY,
    user_id integer NOT NULL,
    body text
);
