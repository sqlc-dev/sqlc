CREATE TABLE employees (
  id       BIGSERIAL                       PRIMARY KEY,
  name     text                            UNIQUE NOT NULL,
  manager  text REFERENCES employees(name)
);