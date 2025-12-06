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
    name String,
    location String
)
ENGINE = MergeTree()
ORDER BY id;

CREATE TABLE IF NOT EXISTS orders
(
    id UInt32,
    user_id UInt32,
    amount Float64,
    created_at DateTime
)
ENGINE = MergeTree()
ORDER BY (id, user_id);
