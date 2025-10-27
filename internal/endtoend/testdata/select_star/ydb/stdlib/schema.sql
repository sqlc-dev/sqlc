CREATE TABLE users (
    id Serial,
    first_name Utf8 NOT NULL,
    last_name Utf8,
    age Int32 NOT NULL,
    PRIMARY KEY (id)
);
