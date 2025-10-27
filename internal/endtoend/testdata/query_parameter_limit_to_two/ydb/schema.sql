-- Example queries for sqlc
CREATE TABLE authors (
    id BigSerial,
    name Text NOT NULL,
    bio Text,
    country_code Text NOT NULL,
    titles Text,
    PRIMARY KEY (id)
);

CREATE TABLE clients (
    id Int32,
    name Text NOT NULL,
    PRIMARY KEY (id)
);
