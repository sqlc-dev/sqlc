CREATE TABLE things (
  id INTEGER PRIMARY KEY,
  title VARCHAR(255),
  code VARYING CHARACTER(10),
  price DECIMAL(10,5),
  flag BOOLEAN,
  big UNSIGNED BIG INT,
  data BLOB,
  untyped,
  ratio DOUBLE PRECISION,
  weird FOO BAR(3),
  created DATETIME,
  n INT NOT NULL,
  label NCHAR(5)
);

CREATE TABLE strict_things (
  id INTEGER,
  anything ANY,
  body TEXT,
  amount REAL,
  raw BLOB,
  n INT
) STRICT;
