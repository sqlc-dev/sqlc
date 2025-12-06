CREATE TABLE IF NOT EXISTS data
(
    id UInt32,
    amount String,
    quantity String,
    created_date String
)
ENGINE = MergeTree()
ORDER BY id;
