CREATE TABLE foo (id integer not null);
CREATE TABLE bar (id integer not null references foo(id), title text);

