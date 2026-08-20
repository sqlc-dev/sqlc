CREATE TABLE inventory (
    sku        text PRIMARY KEY,
    quantity   integer NOT NULL,
    updated_by text
);

CREATE TABLE inventory_updates (
    sku   text PRIMARY KEY,
    delta integer NOT NULL
);

-- Tables for the RETURNING * column-ordering regression. The target and
-- source share a column name ("id") with distinguishable types so the
-- generated field types prove the source-before-target expansion order.
CREATE TABLE product_catalog (
    id    text PRIMARY KEY,
    name  text NOT NULL,
    price integer NOT NULL
);

CREATE TABLE product_imports (
    id    integer NOT NULL,
    label text NOT NULL
);
