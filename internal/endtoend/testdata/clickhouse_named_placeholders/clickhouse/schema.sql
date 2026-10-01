CREATE TABLE events (
    id UInt64,
    name String,
    tag Nullable(String),
    amount Decimal(10, 2),
    uid UUID,
    ip IPv4,
    big Int128
) ENGINE = MergeTree ORDER BY id;
