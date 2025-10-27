-- Example queries for sqlc
CREATE TABLE authors (
    id BigSerial,
    name Text NOT NULL,
    bio Text,
    country_code Text NOT NULL,
    PRIMARY KEY (id)
);
