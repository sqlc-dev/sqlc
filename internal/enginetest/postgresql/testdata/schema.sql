-- =============================================================================
-- PostgreSQL Core Schema
-- This schema is used by all PostgreSQL end-to-end tests
-- =============================================================================

-- Users table: core entity with various column types
CREATE TABLE users (
    id              SERIAL PRIMARY KEY,
    username        VARCHAR(255) NOT NULL UNIQUE,
    email           VARCHAR(255) NOT NULL,
    full_name       TEXT,
    age             INT,
    balance         DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    role            VARCHAR(20) NOT NULL DEFAULT 'user',
    bio             TEXT,
    metadata        BYTEA,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP,
    deleted_at      TIMESTAMP
);

-- Categories table: self-referential for tree structures
CREATE TABLE categories (
    id              SERIAL PRIMARY KEY,
    parent_id       INT REFERENCES categories(id),
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    sort_order      INT NOT NULL DEFAULT 0,
    is_visible      BOOLEAN NOT NULL DEFAULT true
);

-- Products table: many-to-one with categories
CREATE TABLE products (
    id              SERIAL PRIMARY KEY,
    category_id     INT REFERENCES categories(id),
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    price           DECIMAL(10,2) NOT NULL,
    quantity        INT NOT NULL DEFAULT 0,
    weight          REAL,
    is_available    BOOLEAN NOT NULL DEFAULT true,
    tags            TEXT,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP
);

-- Orders table: many-to-one with users
CREATE TABLE orders (
    id              SERIAL PRIMARY KEY,
    user_id         INT NOT NULL REFERENCES users(id),
    status          VARCHAR(20) NOT NULL DEFAULT 'draft',
    total_amount    DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    notes           TEXT,
    shipped_at      TIMESTAMP,
    delivered_at    TIMESTAMP,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP
);

-- Order items table: many-to-one with orders and products
CREATE TABLE order_items (
    id              SERIAL PRIMARY KEY,
    order_id        INT NOT NULL REFERENCES orders(id),
    product_id      INT NOT NULL REFERENCES products(id),
    quantity        INT NOT NULL DEFAULT 1,
    unit_price      DECIMAL(10,2) NOT NULL,
    discount        DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    UNIQUE(order_id, product_id)
);

-- Tags table: for many-to-many relationships
CREATE TABLE tags (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL UNIQUE,
    color           VARCHAR(7),
    created_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Product tags junction table: many-to-many
CREATE TABLE product_tags (
    product_id      INT NOT NULL REFERENCES products(id),
    tag_id          INT NOT NULL REFERENCES tags(id),
    PRIMARY KEY (product_id, tag_id)
);

-- Audit log table
CREATE TABLE audit_logs (
    id              SERIAL PRIMARY KEY,
    table_name      VARCHAR(100) NOT NULL,
    record_id       INT NOT NULL,
    action          VARCHAR(20) NOT NULL,
    old_values      TEXT,
    new_values      TEXT,
    user_id         INT REFERENCES users(id),
    ip_address      VARCHAR(45),
    created_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Settings table: key-value store
CREATE TABLE settings (
    key             VARCHAR(255) PRIMARY KEY,
    value           TEXT NOT NULL,
    value_type      VARCHAR(20) NOT NULL DEFAULT 'string',
    description     TEXT,
    updated_at      TIMESTAMP
);

-- Tasks table
CREATE TABLE tasks (
    id              SERIAL PRIMARY KEY,
    user_id         INT NOT NULL REFERENCES users(id),
    title           VARCHAR(255) NOT NULL,
    description     TEXT,
    priority        VARCHAR(20) NOT NULL DEFAULT 'medium',
    is_completed    BOOLEAN NOT NULL DEFAULT false,
    due_date        DATE,
    started_at      TIMESTAMP,
    completed_at    TIMESTAMP,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

-- =============================================================================
-- Views
-- =============================================================================

CREATE VIEW active_users AS
    SELECT id, username, email, full_name, created_at
    FROM users
    WHERE is_active = true AND deleted_at IS NULL;

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
