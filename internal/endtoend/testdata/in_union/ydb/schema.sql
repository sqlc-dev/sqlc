CREATE TABLE authors (
    id Int32,
    name Text NOT NULL,
    bio Text,
    PRIMARY KEY (id)
);

CREATE TABLE book1 (
    author_id Int32,
    name Text,
    PRIMARY KEY (author_id)
);

CREATE TABLE book2 (
    author_id Int32,
    name Text,
    PRIMARY KEY (author_id)
);


