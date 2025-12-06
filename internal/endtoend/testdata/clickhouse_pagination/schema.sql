CREATE TABLE IF NOT EXISTS products
(
    id UInt32,
    name String,
    price Float32,
    stock UInt32
)
ENGINE = MergeTree()
ORDER BY id;
