-- UUID Types
-- https://ydb.tech/docs/ru/concepts/datatypes#uuid

CREATE TABLE dt_uuid (
    id Serial,
    a Uuid,
    PRIMARY KEY (id)
);

CREATE TABLE dt_uuid_not_null (
    id Serial,
    a Uuid NOT NULL,
    PRIMARY KEY (id)
);
