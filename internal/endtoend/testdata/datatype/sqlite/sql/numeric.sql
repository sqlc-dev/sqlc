-- Numeric Types
-- https://www.sqlite.org/datatype3.html
CREATE TABLE dt_numeric (
    a INT,
    b INTEGER,
    c TINYINT,
    d SMALLINT,
    e MEDIUMINT,
    f BIGINT,
    g UNSIGNED BIG INT,
    h INT2,
    i INT8,
    j REAL,
    k DOUBLE,
    l DOUBLE PRECISION,
    m FLOAT,
    n NUMERIC,
    o DECIMAL(10,5)
);

CREATE TABLE dt_numeric_not_null (
    a INT NOT NULL,
    b INTEGER NOT NULL,
    c TINYINT NOT NULL,
    d SMALLINT NOT NULL,
    e MEDIUMINT NOT NULL,
    f BIGINT NOT NULL,
    g UNSIGNED BIG INT NOT NULL,
    h INT2 NOT NULL,
    i INT8 NOT NULL,
    j REAL NOT NULL,
    k DOUBLE NOT NULL,
    l DOUBLE PRECISION NOT NULL,
    m FLOAT NOT NULL,
    n NUMERIC NOT NULL,
    o DECIMAL(10,5) NOT NULL
);
