-- JSON and Document Types
-- https://ydb.tech/docs/ru/concepts/datatypes#json

CREATE TABLE dt_json (
    id Serial,
    a Json,
    b JsonDocument,
    c Yson,
    PRIMARY KEY (id)
);

CREATE TABLE dt_json_not_null (
    id Serial,
    a Json NOT NULL,
    b JsonDocument NOT NULL,
    c Yson NOT NULL,
    PRIMARY KEY (id)
);
