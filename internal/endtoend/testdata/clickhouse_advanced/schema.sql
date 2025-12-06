CREATE TABLE IF NOT EXISTS transactions
(
    id UInt32,
    user_id UInt32,
    amount Float64,
    category String,
    created_at DateTime
)
ENGINE = MergeTree()
ORDER BY (id, created_at);

CREATE TABLE IF NOT EXISTS activities
(
    id UInt32,
    user_id UInt32,
    action String,
    timestamp DateTime
)
ENGINE = MergeTree()
ORDER BY (id, timestamp);
