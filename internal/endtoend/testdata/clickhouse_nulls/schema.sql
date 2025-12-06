CREATE TABLE IF NOT EXISTS products
(
    id UInt32,
    name String,
    description Nullable(String),
    discount Nullable(Float32),
    category String
)
ENGINE = MergeTree()
ORDER BY id;
