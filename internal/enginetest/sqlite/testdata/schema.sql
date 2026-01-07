-- =============================================================================
-- SQLite Core Schema
-- This schema is used by all SQLite end-to-end tests
-- =============================================================================

-- Users table: core entity with various column types
CREATE TABLE users (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    username        TEXT NOT NULL UNIQUE,
    email           TEXT NOT NULL,
    full_name       TEXT,
    age             INTEGER,
    balance         REAL NOT NULL DEFAULT 0.00,
    is_active       INTEGER NOT NULL DEFAULT 1,
    status          TEXT NOT NULL DEFAULT 'pending',
    role            TEXT NOT NULL DEFAULT 'user',
    bio             TEXT,
    metadata        BLOB,
    created_at      TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at      TEXT,
    deleted_at      TEXT
);

-- Categories table: self-referential for tree structures
CREATE TABLE categories (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    parent_id       INTEGER REFERENCES categories(id),
    name            TEXT NOT NULL,
    description     TEXT,
    sort_order      INTEGER NOT NULL DEFAULT 0,
    is_visible      INTEGER NOT NULL DEFAULT 1
);

-- Products table: many-to-one with categories
CREATE TABLE products (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    category_id     INTEGER REFERENCES categories(id),
    name            TEXT NOT NULL,
    description     TEXT,
    price           REAL NOT NULL,
    quantity        INTEGER NOT NULL DEFAULT 0,
    weight          REAL,
    is_available    INTEGER NOT NULL DEFAULT 1,
    tags            TEXT,
    created_at      TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at      TEXT
);

-- Orders table: many-to-one with users
CREATE TABLE orders (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id         INTEGER NOT NULL REFERENCES users(id),
    status          TEXT NOT NULL DEFAULT 'draft',
    total_amount    REAL NOT NULL DEFAULT 0.00,
    notes           TEXT,
    shipped_at      TEXT,
    delivered_at    TEXT,
    created_at      TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at      TEXT
);

-- Order items table: many-to-one with orders and products
CREATE TABLE order_items (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    order_id        INTEGER NOT NULL REFERENCES orders(id),
    product_id      INTEGER NOT NULL REFERENCES products(id),
    quantity        INTEGER NOT NULL DEFAULT 1,
    unit_price      REAL NOT NULL,
    discount        REAL NOT NULL DEFAULT 0.00,
    UNIQUE(order_id, product_id)
);

-- Tags table: for many-to-many relationships
CREATE TABLE tags (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    name            TEXT NOT NULL UNIQUE,
    color           TEXT,
    created_at      TEXT NOT NULL DEFAULT (datetime('now'))
);

-- Product tags junction table: many-to-many
CREATE TABLE product_tags (
    product_id      INTEGER NOT NULL REFERENCES products(id),
    tag_id          INTEGER NOT NULL REFERENCES tags(id),
    PRIMARY KEY (product_id, tag_id)
);

-- Audit log table
CREATE TABLE audit_logs (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    table_name      TEXT NOT NULL,
    record_id       INTEGER NOT NULL,
    action          TEXT NOT NULL,
    old_values      TEXT,
    new_values      TEXT,
    user_id         INTEGER REFERENCES users(id),
    ip_address      TEXT,
    created_at      TEXT NOT NULL DEFAULT (datetime('now'))
);

-- Settings table: key-value store
CREATE TABLE settings (
    key             TEXT PRIMARY KEY,
    value           TEXT NOT NULL,
    value_type      TEXT NOT NULL DEFAULT 'string',
    description     TEXT,
    updated_at      TEXT
);

-- Tasks table
CREATE TABLE tasks (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id         INTEGER NOT NULL REFERENCES users(id),
    title           TEXT NOT NULL,
    description     TEXT,
    priority        TEXT NOT NULL DEFAULT 'medium',
    is_completed    INTEGER NOT NULL DEFAULT 0,
    due_date        TEXT,
    started_at      TEXT,
    completed_at    TEXT,
    created_at      TEXT NOT NULL DEFAULT (datetime('now'))
);

-- =============================================================================
-- Views
-- =============================================================================

CREATE VIEW active_users AS
    SELECT id, username, email, full_name, created_at
    FROM users
    WHERE is_active = 1 AND deleted_at IS NULL;

CREATE VIEW order_summaries AS
    SELECT
        o.id AS order_id,
        o.user_id,
        u.username,
        o.status,
        o.total_amount,
        COUNT(oi.id) AS item_count,
        o.created_at
    FROM orders o
    JOIN users u ON o.user_id = u.id
    LEFT JOIN order_items oi ON o.id = oi.order_id
    GROUP BY o.id, o.user_id, u.username, o.status, o.total_amount, o.created_at;

CREATE VIEW category_tree AS
    SELECT
        c.id,
        c.name,
        c.parent_id,
        p.name AS parent_name
    FROM categories c
    LEFT JOIN categories p ON c.parent_id = p.id;
