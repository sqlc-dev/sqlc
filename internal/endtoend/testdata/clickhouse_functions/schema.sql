CREATE TABLE IF NOT EXISTS web_logs
(
    id UInt32,
    ip_address IPv4,
    timestamp DateTime,
    url String,
    response_time_ms UInt32,
    status_code UInt16,
    user_agent String
)
ENGINE = MergeTree()
ORDER BY (timestamp, ip_address);

CREATE TABLE IF NOT EXISTS metrics
(
    name String,
    value Float64,
    tags Map(String, String),
    created_at DateTime
)
ENGINE = MergeTree()
ORDER BY (created_at, name);
