CREATE TABLE authors (
  id   BigSerial,
  name Text NOT NULL,
  bio  Text,
  PRIMARY KEY (id)
);

CREATE TABLE books (
  id BigSerial,
  author_id BigSerial NOT NULL,
  title Text NOT NULL,
  PRIMARY KEY (id)
);



