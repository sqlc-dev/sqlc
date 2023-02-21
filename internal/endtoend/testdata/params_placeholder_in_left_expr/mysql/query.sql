CREATE TABLE users (
    id INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255)
) ENGINE=InnoDB;

-- name: FindByID :many
SELECT * FROM users WHERE ? = id;

-- name: FindByIDAndName :many
SELECT * FROM users WHERE ? = id AND ? = name;
