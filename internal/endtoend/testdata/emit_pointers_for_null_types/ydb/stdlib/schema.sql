-- Numeric Types for YDB
CREATE TABLE dt_numeric (
    a Int8,
    b Int16,
    c Int32,
    d Int64,
    e Uint8,
    f Uint16,
    g Uint32,
    h Uint64,
    i Float,
    j Double,
    k Decimal(10, 2),
    l SmallSerial,
    m Serial,
    n BigSerial,
    PRIMARY KEY (a)
);

CREATE TABLE dt_numeric_not_null (
    a Int8 NOT NULL,
    b Int16 NOT NULL,
    c Int32 NOT NULL,
    d Int64 NOT NULL,
    e Uint8 NOT NULL,
    f Uint16 NOT NULL,
    g Uint32 NOT NULL,
    h Uint64 NOT NULL,
    i Float NOT NULL,
    j Double NOT NULL,
    k Decimal(10, 2) NOT NULL,
    l SmallSerial NOT NULL,
    m Serial NOT NULL,
    n BigSerial NOT NULL,
    PRIMARY KEY (a)
);
