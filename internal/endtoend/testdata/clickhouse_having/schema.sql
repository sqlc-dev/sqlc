CREATE TABLE IF NOT EXISTS sales
(
    id UInt32,
    region String,
    amount Float32,
    year UInt32
)
ENGINE = MergeTree()
ORDER BY id;
