-- Numeric Types
-- https://ydb.tech/docs/ru/concepts/datatypes#numeric

CREATE TABLE dt_numeric (
    id Serial,
    -- Integer types
    a Int8,
    b Int16,
    c Int32,
    d Int64,
    e Uint8,
    f Uint16,
    g Uint32,
    h Uint64,
    -- Floating point types
    i Float,
    j Double,
    -- Decimal types
    k Decimal(22, 9),
    l Decimal(35, 0),
    -- Serial types
    n SmallSerial,
    o Serial2,
    p Serial,
    q Serial4,
    r Serial8,
    s BigSerial,
    PRIMARY KEY (id)
);

CREATE TABLE dt_numeric_not_null (
    id Serial,
    -- Integer types
    a Int8 NOT NULL,
    b Int16 NOT NULL,
    c Int32 NOT NULL,
    d Int64 NOT NULL,
    e Uint8 NOT NULL,
    f Uint16 NOT NULL,
    g Uint32 NOT NULL,
    h Uint64 NOT NULL,
    -- Floating point types
    i Float NOT NULL,
    j Double NOT NULL,
    -- Decimal types
    k Decimal(22, 9) NOT NULL,
    l Decimal(35, 0) NOT NULL,
    -- Serial types
    n SmallSerial NOT NULL,
    o Serial2 NOT NULL,
    p Serial NOT NULL,
    q Serial4 NOT NULL,
    r Serial8 NOT NULL,
    s BigSerial NOT NULL,
    PRIMARY KEY (id)
);

