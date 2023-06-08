-- Range Types
-- https://www.postgresql.org/docs/current/rangetypes.html
CREATE TABLE dt_range (
    a int4range,
    b int8range,
    c numrange,
    d tsrange,
    e tstzrange,
    f daterange
);

CREATE TABLE dt_range_not_null (
    a int4range NOT NULL,
    b int8range NOT NULL,
    c numrange NOT NULL,
    d tsrange NOT NULL,
    e tstzrange NOT NULL,
    f daterange NOT NULL
);
