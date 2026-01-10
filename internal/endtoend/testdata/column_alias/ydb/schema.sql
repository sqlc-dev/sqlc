CREATE TABLE users (
  id Serial,
  fname Text NOT NULL,
  lname Text NOT NULL,
  email Text NOT NULL,
  enc_passwd Text NOT NULL,
  created_at DateTime NOT NULL,
  PRIMARY KEY (id)
);
