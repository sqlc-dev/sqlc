CREATE TABLE foo (id integer not null, bar_id integer references bar(id));
CREATE TABLE bar (id integer not null);

