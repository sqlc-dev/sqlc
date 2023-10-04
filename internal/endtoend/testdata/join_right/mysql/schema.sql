CREATE TABLE foo (id serial not null, bar_id int references bar(id));
CREATE TABLE bar (id serial not null);

