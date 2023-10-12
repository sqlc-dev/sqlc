CREATE TABLE orders (
  id   BIGSERIAL PRIMARY KEY,
  quantity decimal      NOT NULL,
  order_catalog  int
);