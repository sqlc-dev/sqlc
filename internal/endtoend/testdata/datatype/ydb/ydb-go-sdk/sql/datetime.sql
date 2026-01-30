-- Date and Time Types
-- https://ydb.tech/docs/ru/concepts/datatypes#datetime

CREATE TABLE dt_datetime (
    id Serial,
    a Date,
    b Date32,
    c Datetime,
    d Datetime64,
    e Timestamp,
    f Timestamp64,
    g Interval,
    h Interval64,
    -- Timezone types
    i TzDate32,
    j TzDateTime64,
    k TzTimestamp64,
    PRIMARY KEY (id)
);

CREATE TABLE dt_datetime_not_null (
    id Serial,
    a Date NOT NULL,
    b Date32 NOT NULL,
    c Datetime NOT NULL,
    d Datetime64 NOT NULL,
    e Timestamp NOT NULL,
    f Timestamp64 NOT NULL,
    g Interval NOT NULL,
    h Interval64 NOT NULL,
    -- Timezone types
    i TzDate32 NOT NULL,
    j TzDateTime64 NOT NULL,
    k TzTimestamp64 NOT NULL,
    PRIMARY KEY (id)
);
