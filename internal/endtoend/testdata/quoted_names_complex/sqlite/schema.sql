-- Test complex quoted table and column names with special characters
-- Covers spaces, hyphens, uppercase, and mixed operations
CREATE TABLE "user profiles" (id integer primary key, data text);
CREATE TABLE "ORDERS" (id integer primary key, data text);
CREATE TABLE products (id integer primary key, data text);
CREATE TABLE "item-categories" (id integer primary key, data text);

-- Test ALTER statements with complex identifiers
ALTER TABLE "user profiles" RENAME COLUMN data TO "profile data";
ALTER TABLE "ORDERS" RENAME TO "customer_orders";
ALTER TABLE products ADD COLUMN "Price Info" text;

-- Test mixed case operations across different statement types
INSERT INTO "user profiles" ("profile data") VALUES ('test data');
UPDATE "ORDERS" SET data = 'updated' WHERE id = 1;
DELETE FROM products WHERE id = 1;

-- Test DROP with various identifier formats
DROP TABLE "user profiles";
DROP TABLE "customer_orders";
DROP TABLE "item-categories";
DROP TABLE products;
