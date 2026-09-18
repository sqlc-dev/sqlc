CREATE SEQUENCE authors_id_seq;

CREATE TABLE authors (
  id   BIGINT PRIMARY KEY DEFAULT nextval('authors_id_seq'),
  name TEXT   NOT NULL,
  bio  TEXT
);
