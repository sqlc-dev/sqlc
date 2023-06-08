-- Numeric Types
-- https://www.postgresql.org/docs/current/datatype-numeric.html
CREATE TABLE dt_numeric (
    -- TODO: this maps incorrectly to int16, not NullInt16
    a smallint,
    b integer,
    c bigint,
    d decimal,
    e numeric,
    f real,
    g double precision,
    -- TODO: this maps incorrectly to int16, not NullInt16
    h smallserial,
    i serial,
    j bigserial,
    -- TODO: this maps incorrectly to int16, not NullInt16
    k int2,
    l int4,
    m int8
);

CREATE TABLE dt_numeric_not_null (
    a smallint NOT NULL,
    b integer NOT NULL,
    c bigint NOT NULL,
    d decimal NOT NULL,
    e numeric NOT NULL,
    f real NOT NULL,
    g double precision NOT NULL,
    h smallserial NOT NULL,
    i serial NOT NULL,
    j bigserial NOT NULL,
    k int2 NOT NULL,
    l int4 NOT NULL,
    m int8 NOT NULL
);
