CREATE TABLE authors (
  id        INT PRIMARY KEY,
  name      TEXT NOT NULL,
  parent_id INT -- nullable
);

