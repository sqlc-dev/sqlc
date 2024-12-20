CREATE TABLE customers (
  id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  first_name varchar(255) not null,
  last_name varchar(255) not null
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;