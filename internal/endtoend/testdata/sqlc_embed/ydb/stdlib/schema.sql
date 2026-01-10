CREATE TABLE users (
    id Int32 NOT NULL,
    name Text NOT NULL,
    age Int32,
    PRIMARY KEY (id)
);

CREATE TABLE posts (
    id Int32 NOT NULL,
    user_id Int32 NOT NULL,
    likes Text NOT NULL,
    PRIMARY KEY (id)
);
