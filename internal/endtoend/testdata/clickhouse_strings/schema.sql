CREATE TABLE IF NOT EXISTS employees
(
    id UInt32,
    first_name String,
    last_name String,
    email String,
    bio String
)
ENGINE = MergeTree()
ORDER BY id;
