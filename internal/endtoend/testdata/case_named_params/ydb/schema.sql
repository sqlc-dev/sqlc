-- https://github.com/sqlc-dev/sqlc/issues/1195

CREATE TABLE authors (
  id   BigSerial,
  username Text,
  email Text,
  name Text NOT NULL,
  bio  Text,
  PRIMARY KEY (id)
);
