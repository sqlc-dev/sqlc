CREATE TABLE IF NOT EXISTS users
(
    id UInt32,
    name String
)
ENGINE = MergeTree()
ORDER BY id;

CREATE TABLE IF NOT EXISTS profiles
(
    user_id UInt32,
    bio String
)
ENGINE = MergeTree()
ORDER BY user_id;
