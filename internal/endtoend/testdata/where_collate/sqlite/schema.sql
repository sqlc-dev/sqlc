CREATE TABLE accounts (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,

    UNIQUE (name COLLATE NOCASE)
);

