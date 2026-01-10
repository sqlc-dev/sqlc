CREATE TABLE foo (barid Serial, PRIMARY KEY (barid));
CREATE TABLE bar (id Serial, owner Text NOT NULL, PRIMARY KEY (id));
