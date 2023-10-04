CREATE TABLE foo (id serial not null);
CREATE TABLE bar (id serial not null references foo(id), title text);

