-- Character Types
-- https://www.sqlite.org/datatype3.html
CREATE TABLE dt_character (
    a CHARACTER(32),
    b VARCHAR(32),
    c VARYING CHARACTER(32),
    d NCHAR(32),
    e NATIVE CHARACTER(32),
    f NVARCHAR(32),
    g TEXT,
    h CLOB
);

CREATE TABLE dt_character_not_null (
    a CHARACTER(32) NOT NULL,
    b VARCHAR(32) NOT NULL,
    c VARYING CHARACTER(32) NOT NULL,
    d NCHAR(32) NOT NULL,
    e NATIVE CHARACTER(32) NOT NULL,
    f NVARCHAR(32) NOT NULL,
    g TEXT NOT NULL,
    h CLOB NOT NULL
);
