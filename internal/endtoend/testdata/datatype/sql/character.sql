-- Character Types
-- https://www.postgresql.org/docs/current/datatype-character.html
CREATE TABLE dt_character (
    a text,
    b character varying(32),
    c varchar(32),
    d character(32),
    e char(32)
);

CREATE TABLE dt_character_not_null (
    a text NOT NULL,
    b character varying(32) NOT NULL,
    c varchar(32) NOT NULL,
    d character(32) NOT NULL,
    e char(32) NOT NULL
);
