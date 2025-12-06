CREATE TABLE IF NOT EXISTS sales
(
    id UInt32,
    product_id UInt32,
    category String,
    amount Float64,
    quantity UInt32,
    created_at DateTime
)
ENGINE = MergeTree()
ORDER BY (id, created_at);
