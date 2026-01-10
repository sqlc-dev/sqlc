-- Character and String Types
-- https://ydb.tech/docs/ru/concepts/datatypes#string

CREATE TABLE dt_character (
    id Serial,
    a String,
    b Utf8,
    PRIMARY KEY (id)
);

CREATE TABLE dt_character_not_null (
    id Serial,
    a String NOT NULL,
    b Utf8 NOT NULL,
    PRIMARY KEY (id)
);

