-- Numeric Types
-- https://dev.mysql.com/doc/refman/8.0/en/numeric-type-syntax.html
CREATE TABLE dt_numeric (
    a INT,
    b INTEGER,
    c TINYINT,
    d SMALLINT,
    e MEDIUMINT,
    f BIGINT,
    g BIT,
    h DECIMAL(10, 5),
    i DEC(10, 5),
    j FLOAT,
    k DOUBLE,
    l DOUBLE PRECISION
);

CREATE TABLE dt_numeric_unsigned (
    a INT UNSIGNED,
    b INTEGER UNSIGNED,
    c TINYINT UNSIGNED,
    d SMALLINT UNSIGNED,
    e MEDIUMINT UNSIGNED,
    f BIGINT UNSIGNED
);

CREATE TABLE dt_numeric_not_null (
    a INT NOT NULL,
    b INTEGER NOT NULL,
    c TINYINT NOT NULL,
    d SMALLINT NOT NULL,
    e MEDIUMINT NOT NULL,
    f BIGINT NOT NULL,
    g BIT NOT NULL,
    h DECIMAL(10, 5) NOT NULL,
    i DEC(10, 5) NOT NULL,
    j FLOAT NOT NULL,
    k DOUBLE NOT NULL,
    l DOUBLE PRECISION NOT NULL
);
