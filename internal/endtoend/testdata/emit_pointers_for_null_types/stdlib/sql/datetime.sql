-- Date/Time Types
-- https://www.postgresql.org/docs/current/datatype-datetime.html
CREATE TABLE dt_datetime (
    a DATE,
    b TIME,
    c TIME WITHOUT TIME ZONE,
    d TIME WITH TIME ZONE,
    e TIMESTAMP,
    f TIMESTAMP WITHOUT TIME ZONE,
    g TIMESTAMP WITH TIME ZONE,
    h timestamptz
);

CREATE TABLE dt_datetime_not_null (
    a DATE NOT NULL,
    b TIME NOT NULL,
    c TIME WITHOUT TIME ZONE NOT NULL,
    d TIME WITH TIME ZONE NOT NULL,
    e TIMESTAMP NOT NULL,
    f TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    g TIMESTAMP WITH TIME ZONE NOT NULL,
    h timestamptz NOT NULL
);
