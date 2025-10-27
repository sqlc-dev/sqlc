-- Example queries for sqlc
CREATE TABLE authors (
  id   Int32,
  name Text NOT NULL,
  age  Int32,
  PRIMARY KEY (id)
);

CREATE TABLE translators (
  id   Int32,
  name Text NOT NULL,
  age  Int32,
  PRIMARY KEY (id)
);

CREATE TABLE books (
  id   Int32,
  author Text NOT NULL,
  translator Text NOT NULL,
  year  Int32,
  PRIMARY KEY (id)
);



