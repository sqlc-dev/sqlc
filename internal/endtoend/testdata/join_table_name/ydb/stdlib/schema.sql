CREATE TABLE bar (id Serial, PRIMARY KEY (id));
CREATE TABLE foo (id Serial, bar Serial, PRIMARY KEY (id));
