CREATE TABLE foo (
  id serial primary key,
  bar text not null,
  baz int not null default 0
);
