CREATE TABLE bar (id serial not null);
CREATE TABLE foo (id serial not null, bar integer references bar(id));

