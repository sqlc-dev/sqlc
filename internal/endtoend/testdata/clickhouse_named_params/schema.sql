CREATE TABLE users
(
    id UInt32,
    name String,
    email String
)
ENGINE = MergeTree()
ORDER BY id;
