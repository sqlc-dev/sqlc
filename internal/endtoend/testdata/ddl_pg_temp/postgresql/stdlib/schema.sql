CREATE TABLE old_table (val SERIAL);
CREATE TABLE new_table (val SERIAL);
CREATE TABLE pg_temp.migrate (val SERIAL);
INSERT INTO pg_temp.migrate (val) SELECT val FROM old_table;
INSERT INTO new_table (val) SELECT val FROM pg_temp.migrate;
