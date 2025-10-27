CREATE TABLE authors (
    id BigSerial,
    name Text NOT NULL,
    bio Text,
    PRIMARY KEY (id)
);

CREATE TABLE people (
    first_name Text NOT NULL,
    PRIMARY KEY (first_name)
);

