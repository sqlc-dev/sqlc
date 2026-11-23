CREATE TABLE authors (
  id         BIGINT PRIMARY KEY AUTO_INCREMENT,
  name       VARCHAR(255) NOT NULL,
  deleted_at DATETIME
);

CREATE TABLE books (
  id         BIGINT PRIMARY KEY AUTO_INCREMENT,
  author_id  BIGINT NOT NULL,
  title      VARCHAR(255) NOT NULL,
  deleted_at DATETIME
);
