CREATE TABLE users (
    id Int64,
    first_name Utf8,
    last_name Utf8,
    age Int64,
    job_status Utf8,
    PRIMARY KEY (id)
);

CREATE TABLE orders (
    id Int64,
    price Decimal(13, 4),
    user_id Int64,
    PRIMARY KEY (id)
);













