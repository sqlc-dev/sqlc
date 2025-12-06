CREATE TABLE IF NOT EXISTS products
(
    id UInt32,
    name String,
    tags Array(String),
    ratings Array(UInt8)
)
ENGINE = MergeTree()
ORDER BY id;

CREATE TABLE IF NOT EXISTS events
(
    id UInt32,
    name String,
    timestamp DateTime,
    properties Nested(
        key String,
        value String
    )
)
ENGINE = MergeTree()
ORDER BY (id, timestamp);
