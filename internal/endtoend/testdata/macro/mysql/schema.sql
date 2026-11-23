CREATE TABLE authors (
  id         BIGINT PRIMARY KEY AUTO_INCREMENT,
  org_id     BIGINT NOT NULL,
  name       VARCHAR(255) NOT NULL,
  bio        TEXT,
  created_at DATETIME NOT NULL,
  deleted_at DATETIME
);
