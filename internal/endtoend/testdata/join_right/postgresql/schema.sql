CREATE TABLE bar (id serial not null unique);
CREATE TABLE foo (id serial not null, bar_id int references bar(id));

