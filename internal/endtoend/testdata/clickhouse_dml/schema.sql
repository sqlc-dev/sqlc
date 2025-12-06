CREATE TABLE IF NOT EXISTS logs
(
    id UInt32,
    level String,
    message String,
    timestamp DateTime,
    source String
)
ENGINE = MergeTree()
ORDER BY (timestamp, id);

CREATE TABLE IF NOT EXISTS notifications
(
    id UInt32,
    user_id UInt32,
    message String,
    read_status UInt8,
    created_at DateTime
)
ENGINE = MergeTree()
ORDER BY (user_id, created_at);
