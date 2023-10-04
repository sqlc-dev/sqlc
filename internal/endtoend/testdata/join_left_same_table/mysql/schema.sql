CREATE TABLE authors (
  id        INT(10) NOT NULL,
  name      VARCHAR(255) NOT NULL,
  parent_id INT(10),
  PRIMARY KEY (id)
);

