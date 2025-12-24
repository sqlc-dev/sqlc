CREATE TABLE authors (
    id UInt64,
    name String,
    bio Nullable(String)
) ENGINE = MergeTree()
ORDER BY id;
