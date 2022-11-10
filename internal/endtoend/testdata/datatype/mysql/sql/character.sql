-- Character Types
-- https://dev.mysql.com/doc/refman/8.0/en/string-type-syntax.html
CREATE TABLE dt_character (
    a CHARACTER(32),
    b VARCHAR(32),
    c CHAR(32),
    d BINARY(32),
    e VARBINARY(32),
    f TINYBLOB,
    g TINYTEXT,
    h TEXT,
    i MEDIUMTEXT,
    j MEDIUMBLOB,
    k LONGTEXT,
    l LONGBLOB
);

CREATE TABLE dt_character_not_null (
    a CHARACTER(32) NOT NULL,
    b VARCHAR(32) NOT NULL,
    c CHAR(32) NOT NULL,
    d BINARY(32) NOT NULL,
    e VARBINARY(32) NOT NULL,
    f TINYBLOB NOT NULL,
    g TINYTEXT NOT NULL,
    h TEXT NOT NULL,
    i MEDIUMTEXT NOT NULL,
    j MEDIUMBLOB NOT NULL,
    k LONGTEXT NOT NULL,
    l LONGBLOB NOT NULL
);
