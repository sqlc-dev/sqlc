CREATE TABLE foo (id serial not null unique);
CREATE TABLE bar (id serial not null references foo(id), title text);

