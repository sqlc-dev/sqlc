CREATE TABLE events (
    id UInt64,
    name String,
    tag Nullable(String),
    amount Float64
) ENGINE = MergeTree ORDER BY id;

CREATE TABLE users (
    id UInt64,
    email String
) ENGINE = MergeTree ORDER BY id;
