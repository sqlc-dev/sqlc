CREATE TABLE old (val SERIAL);
CREATE TABLE new (val SERIAL);
CREATE TABLE pg_temp.migrate (val SERIAL);
INSERT INTO pg_temp.migrate (val) SELECT val FROM old;
INSERT INTO new (val) SELECT val FROM pg_temp.migrate;
