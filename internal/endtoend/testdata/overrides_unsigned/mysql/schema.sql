CREATE TABLE authors (
  id     SERIAL,
  name   text NOT NULL,
  rating bigint NOT NULL,
  score  int UNSIGNED NOT NULL,
  bio  text
);

