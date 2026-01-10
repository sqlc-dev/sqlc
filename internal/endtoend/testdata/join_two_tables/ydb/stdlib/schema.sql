CREATE TABLE foo (bar_id Serial, baz_id Serial, PRIMARY KEY (bar_id, baz_id));
CREATE TABLE bar (id Serial, PRIMARY KEY (id));
CREATE TABLE baz (id Serial, PRIMARY KEY (id));
