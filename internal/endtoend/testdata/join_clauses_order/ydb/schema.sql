CREATE TABLE a (
  id BigSerial,
  a Text,
  PRIMARY KEY (id)
);

CREATE TABLE b (
  id BigSerial,
  b Text,
  a_id Int64,
  PRIMARY KEY (id)
);

CREATE TABLE c (
  id BigSerial,
  c Text,
  a_id Int64,
  PRIMARY KEY (id)
);

CREATE TABLE d (
  id BigSerial,
  d Text,
  a_id Int64,
  PRIMARY KEY (id)
);

CREATE TABLE e (
  id BigSerial,
  e Text,
  a_id Int64,
  PRIMARY KEY (id)
);
