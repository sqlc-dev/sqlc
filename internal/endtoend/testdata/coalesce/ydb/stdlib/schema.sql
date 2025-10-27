CREATE TABLE foo (
  bar Text,
  bat Text NOT NULL,
  baz Int64,
  qux Int64 NOT NULL,
  PRIMARY KEY (bat, qux)
);
