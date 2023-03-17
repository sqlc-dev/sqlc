CREATE TABLE authors (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name text NOT NULL,
  deleted_at datetime NOT NULL,
  updated_at datetime NOT NULL
);

CREATE TABLE books (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  is_amazing tinyint(1) NOT NULL
);