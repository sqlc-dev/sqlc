CREATE TABLE IF NOT EXISTS orders
(
    id UInt32,
    customer_name String,
    total Float32,
    created_at DateTime
)
ENGINE = MergeTree()
ORDER BY id;
