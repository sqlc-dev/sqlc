CREATE TABLE IF NOT EXISTS users
(
    id UInt32,
    name String,
    department String,
    salary UInt32
)
ENGINE = MergeTree()
ORDER BY id;
