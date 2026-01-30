-- Boolean Types
-- https://ydb.tech/docs/ru/concepts/datatypes#boolean

CREATE TABLE dt_boolean (
    id Serial,
    a Bool,
    PRIMARY KEY (id)
);

CREATE TABLE dt_boolean_not_null (
    id Serial,
    a Bool NOT NULL,
    PRIMARY KEY (id)
);
