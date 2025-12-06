CREATE TABLE IF NOT EXISTS logs
(
    id UInt32,
    level String,
    message String,
    created_at DateTime
)
ENGINE = MergeTree()
ORDER BY id;
