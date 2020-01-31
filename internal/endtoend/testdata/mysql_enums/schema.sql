CREATE TABLE examples (
    first_name ENUM('john', 'albert') NOT NULL,
    user_id ENUM('one', 'two') NOT NULL,
    last_name ENUM('smith', 'frank') NOT NULL
) ENGINE=InnoDB;
