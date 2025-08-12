-- Test table name case sensitivity handling across different SQLite operations
-- Create tables with different case patterns to verify consistent name resolution
CREATE TABLE users (id integer primary key, name text);
CREATE TABLE "Authors" (id integer primary key, name text);
CREATE TABLE Books (id integer primary key, title text);

-- Create a temporary table to test drop operations  
CREATE TABLE temp_orders (id integer primary key);
DROP TABLE temp_orders;

-- Create another temp table with quoted identifier
CREATE TABLE "temp_products" (id integer primary key);
DROP TABLE "temp_products";
