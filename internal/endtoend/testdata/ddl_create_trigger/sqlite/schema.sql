/* examples copied from https://www.sqlite.org/lang_createtrigger.html
   only expectation in sqlc is that they parse, codegen is unaffected */

CREATE TABLE trigger_customers (
    name TEXT PRIMARY KEY,
    address TEXT
);

CREATE TABLE trigger_orders (
    id INTEGER PRIMARY KEY,
    customer_name TEXT,
    address TEXT
);

CREATE TRIGGER update_customer_address UPDATE OF address ON trigger_customers
BEGIN
    UPDATE trigger_orders SET address = new.address WHERE customer_name = old.name;
END;

CREATE TABLE customer(
                         cust_id INTEGER PRIMARY KEY,
                         cust_name TEXT,
                         cust_addr TEXT
);
CREATE VIEW customer_address AS
SELECT cust_id, cust_addr FROM customer;
CREATE TRIGGER cust_addr_chng
    INSTEAD OF UPDATE OF cust_addr ON customer_address
BEGIN
    UPDATE customer SET cust_addr=NEW.cust_addr
    WHERE cust_id=NEW.cust_id;
END;