-- Date/Time Types
-- https://www.sqlite.org/datatype3.html
CREATE TABLE dt_datetime (
    a DATE,
    b DATETIME,
    c TIMESTAMP
);

CREATE TABLE dt_datetime_not_null (
    a DATE NOT NULL,
    b DATETIME NOT NULL,
    c TIMESTAMP NOT NULL
);
