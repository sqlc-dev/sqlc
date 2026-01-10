CREATE TABLE foo (
  id Int64,
  PRIMARY KEY (id)
);

CREATE TABLE bar (
  foo_id Int64,
  info Text,
  PRIMARY KEY (foo_id)
);
