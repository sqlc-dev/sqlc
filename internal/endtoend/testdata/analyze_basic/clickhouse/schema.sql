CREATE TABLE events (
    id UInt64,
    name String,
    tag Nullable(String),
    amount Float64,
    created DateTime
) ENGINE = MergeTree ORDER BY id;
