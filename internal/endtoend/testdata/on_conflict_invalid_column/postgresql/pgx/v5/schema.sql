CREATE TABLE cart_items (
    owner_id       VARCHAR(255) NOT NULL,
    product_id     UUID         NOT NULL,
    price_amount   DECIMAL      NOT NULL,
    price_currency VARCHAR(3)   NOT NULL,
    PRIMARY KEY (owner_id, product_id)
);
