CREATE TABLE orders (
    order_id SERIAL PRIMARY KEY,
    customer_name TEXT,
    amount DECIMAL(10, 2)
);

CREATE TABLE shipments (
    shipment_id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(order_id),
    address TEXT
);
