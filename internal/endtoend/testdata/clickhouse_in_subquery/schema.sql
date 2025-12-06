CREATE TABLE IF NOT EXISTS users
(
    id UInt32,
    name String,
    department_id UInt32
)
ENGINE = MergeTree()
ORDER BY id;

CREATE TABLE IF NOT EXISTS departments
(
    id UInt32,
    name String
)
ENGINE = MergeTree()
ORDER BY id;
