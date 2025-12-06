CREATE TABLE IF NOT EXISTS customers
(
    id UInt32,
    name String,
    city String
)
ENGINE = MergeTree()
ORDER BY id;

CREATE TABLE IF NOT EXISTS orders
(
    id UInt32,
    customer_id UInt32,
    amount Float32
)
ENGINE = MergeTree()
ORDER BY id;
