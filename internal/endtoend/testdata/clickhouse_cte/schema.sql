CREATE TABLE IF NOT EXISTS employees
(
    id UInt32,
    name String,
    manager_id Nullable(UInt32),
    salary UInt32
)
ENGINE = MergeTree()
ORDER BY id;
