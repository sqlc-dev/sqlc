-- https://github.com/sqlc-dev/sqlc/issues/4117
CREATE TABLE a (
    id uuid PRIMARY KEY,
    name TEXT
);

CREATE TABLE b (
    id uuid PRIMARY KEY,
    a_id uuid NOT NULL REFERENCES a (id),
    name TEXT
);
