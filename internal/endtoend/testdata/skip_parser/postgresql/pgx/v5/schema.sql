CREATE TABLE products (
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT NOT NULL,
    price      NUMERIC(10, 2) NOT NULL,
    tags       TEXT[]
);
