-- =============================================================================
-- MySQL Core Schema
-- This schema is used by all MySQL end-to-end tests
-- =============================================================================

-- Users table: core entity with various column types
CREATE TABLE users (
    id              INT AUTO_INCREMENT PRIMARY KEY,
    username        VARCHAR(255) NOT NULL UNIQUE,
    email           VARCHAR(255) NOT NULL,
    full_name       TEXT,
    age             INT,
    balance         DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    role            VARCHAR(20) NOT NULL DEFAULT 'user',
    bio             TEXT,
    metadata        BLOB,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NULL,
    deleted_at      TIMESTAMP NULL
);

-- Categories table: self-referential for tree structures
CREATE TABLE categories (
    id              INT AUTO_INCREMENT PRIMARY KEY,
    parent_id       INT,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    sort_order      INT NOT NULL DEFAULT 0,
    is_visible      BOOLEAN NOT NULL DEFAULT true,
    FOREIGN KEY (parent_id) REFERENCES categories(id)
);

-- Products table: many-to-one with categories
CREATE TABLE products (
    id              INT AUTO_INCREMENT PRIMARY KEY,
    category_id     INT,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    price           DECIMAL(10,2) NOT NULL,
    quantity        INT NOT NULL DEFAULT 0,
    weight          FLOAT,
    is_available    BOOLEAN NOT NULL DEFAULT true,
    tags            TEXT,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NULL,
    FOREIGN KEY (category_id) REFERENCES categories(id)
);

-- Orders table: many-to-one with users
CREATE TABLE orders (
    id              INT AUTO_INCREMENT PRIMARY KEY,
    user_id         INT NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'draft',
    total_amount    DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    notes           TEXT,
    shipped_at      TIMESTAMP NULL,
    delivered_at    TIMESTAMP NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Order items table: many-to-one with orders and products
CREATE TABLE order_items (
    id              INT AUTO_INCREMENT PRIMARY KEY,
    order_id        INT NOT NULL,
    product_id      INT NOT NULL,
    quantity        INT NOT NULL DEFAULT 1,
    unit_price      DECIMAL(10,2) NOT NULL,
    discount        DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    UNIQUE KEY (order_id, product_id),
    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
);

-- Tags table: for many-to-many relationships
CREATE TABLE tags (
    id              INT AUTO_INCREMENT PRIMARY KEY,
    name            VARCHAR(100) NOT NULL UNIQUE,
    color           VARCHAR(7),
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Product tags junction table: many-to-many
CREATE TABLE product_tags (
    product_id      INT NOT NULL,
    tag_id          INT NOT NULL,
    PRIMARY KEY (product_id, tag_id),
    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (tag_id) REFERENCES tags(id)
);

-- Audit log table
CREATE TABLE audit_logs (
    id              INT AUTO_INCREMENT PRIMARY KEY,
    table_name      VARCHAR(100) NOT NULL,
    record_id       INT NOT NULL,
    action          VARCHAR(20) NOT NULL,
    old_values      TEXT,
    new_values      TEXT,
    user_id         INT,
    ip_address      VARCHAR(45),
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Settings table: key-value store
CREATE TABLE settings (
    `key`           VARCHAR(255) PRIMARY KEY,
    value           TEXT NOT NULL,
    value_type      VARCHAR(20) NOT NULL DEFAULT 'string',
    description     TEXT,
    updated_at      TIMESTAMP NULL
);

-- Tasks table
CREATE TABLE tasks (
    id              INT AUTO_INCREMENT PRIMARY KEY,
    user_id         INT NOT NULL,
    title           VARCHAR(255) NOT NULL,
    description     TEXT,
    priority        VARCHAR(20) NOT NULL DEFAULT 'medium',
    is_completed    BOOLEAN NOT NULL DEFAULT false,
    due_date        DATE,
    started_at      TIMESTAMP NULL,
    completed_at    TIMESTAMP NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
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
