CREATE TABLE bar (id integer not null);
CREATE TABLE foo (id integer not null, bar integer references bar(id));

