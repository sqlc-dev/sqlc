CREATE TABLE inventory (
    sku        text PRIMARY KEY,
    quantity   integer NOT NULL,
    updated_by text
);

CREATE TABLE inventory_updates (
    sku   text PRIMARY KEY,
    delta integer NOT NULL
);
