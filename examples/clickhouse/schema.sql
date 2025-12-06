-- ClickHouse example schema

CREATE DATABASE IF NOT EXISTS sqlc_example;

-- Basic CRUD tables
CREATE TABLE IF NOT EXISTS sqlc_example.users
(
    id UInt32,
    name String,
    email String,
    status String,
    created_at DateTime
)
ENGINE = MergeTree()
ORDER BY id;

CREATE TABLE IF NOT EXISTS sqlc_example.posts
(
    id UInt32,
    user_id UInt32,
    title String,
    content String,
    status Nullable(String),
    created_at DateTime
)
ENGINE = MergeTree()
ORDER BY (id, user_id);

CREATE TABLE IF NOT EXISTS sqlc_example.comments
(
    id UInt32,
    post_id UInt32,
    user_id UInt32,
    content String,
    created_at DateTime
)
ENGINE = MergeTree()
ORDER BY (id, post_id, user_id);

-- Tables with array columns for ARRAY JOIN examples
CREATE TABLE IF NOT EXISTS sqlc_example.users_with_tags
(
    id UInt32,
    name String,
    tags Array(String)
)
ENGINE = MergeTree()
ORDER BY id;

CREATE TABLE IF NOT EXISTS sqlc_example.events_with_properties
(
    event_id UInt32,
    event_name String,
    timestamp DateTime,
    properties Nested(
        keys String,
        values String
    )
)
ENGINE = MergeTree()
ORDER BY event_id;

CREATE TABLE IF NOT EXISTS sqlc_example.nested_table
(
    record_id UInt32,
    nested_array Array(String)
)
ENGINE = MergeTree()
ORDER BY record_id;

CREATE TABLE IF NOT EXISTS sqlc_example.products
(
    product_id UInt32,
    name String,
    categories Array(String)
)
ENGINE = MergeTree()
ORDER BY product_id;

-- Metrics and analytics tables
CREATE TABLE IF NOT EXISTS sqlc_example.metrics
(
    category String,
    value Float64,
    value_x Float64,
    value_y Float64,
    timestamp DateTime
)
ENGINE = MergeTree()
ORDER BY timestamp;



CREATE TABLE IF NOT EXISTS sqlc_example.order_metrics
(
    status String,
    amount Float64,
    rating Nullable(Float64),
    created_at DateTime
)
ENGINE = MergeTree()
ORDER BY created_at;

CREATE TABLE IF NOT EXISTS sqlc_example.timeseries
(
    date Date,
    metric_value Float64
)
ENGINE = MergeTree()
ORDER BY date;

CREATE TABLE IF NOT EXISTS sqlc_example.events
(
    id UInt32,
    amount Float64,
    created_at DateTime,
    status String,
    platform_id Nullable(String),
    json_value Nullable(String)
)
ENGINE = MergeTree()
ORDER BY id;

-- Table with Map type columns
CREATE TABLE IF NOT EXISTS sqlc_example.config_data
(
    id UInt32,
    settings Map(String, String),
    metrics Map(String, Float64),
    tags Map(String, Array(String)),
    created_at DateTime
)
ENGINE = MergeTree()
ORDER BY id;

-- Table with IP address columns
CREATE TABLE IF NOT EXISTS sqlc_example.network_logs
(
    id UInt32,
    source_ip IPv4,
    dest_ip IPv4,
    source_ipv6 Nullable(IPv6),
    dest_ipv6 Nullable(IPv6),
    timestamp DateTime,
    bytes_sent UInt64
)
ENGINE = MergeTree()
ORDER BY (timestamp, source_ip);

-- Event tracking tables for complex UNION/CTE queries
CREATE TABLE IF NOT EXISTS sqlc_example.platform_created_event
(
    event_id UInt32,
    timestamp DateTime,
    platform_id UInt32,
    platform_name String,
    region String
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (timestamp, platform_id);

CREATE TABLE IF NOT EXISTS sqlc_example.platform_renamed_event
(
    event_id UInt32,
    timestamp DateTime,
    platform_id UInt32,
    old_name String,
    new_name String,
    region String
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (timestamp, platform_id);

CREATE TABLE IF NOT EXISTS sqlc_example.feature_usage
(
    event_id UInt32,
    timestamp DateTime,
    platform_id UInt32,
    feature_id String,
    user_count UInt32,
    usage_count UInt64
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (timestamp, platform_id, feature_id);

-- Tables for LEFT JOIN USING clause test
CREATE TABLE IF NOT EXISTS sqlc_example.orders
(
    order_id UInt32,
    customer_name String,
    amount Float64,
    created_at DateTime
)
ENGINE = MergeTree()
ORDER BY order_id;

CREATE TABLE IF NOT EXISTS sqlc_example.shipments
(
    shipment_id UInt32,
    order_id UInt32,
    address String,
    shipped_at DateTime
)
ENGINE = MergeTree()
ORDER BY (order_id, shipment_id);
