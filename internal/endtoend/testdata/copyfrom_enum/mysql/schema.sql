CREATE TABLE records (
    id INT NOT NULL PRIMARY KEY,
    kind ENUM('start_point', 'end_point') NOT NULL,
    optional_kind ENUM('start_point', 'end_point'),
    note TEXT NOT NULL
);
CREATE TABLE single_values (kind ENUM('start_point', 'end_point') NOT NULL);
CREATE TABLE nullable_values (kind ENUM('start_point', 'end_point'));
CREATE TABLE overridden_values (id INT NOT NULL, kind ENUM('start_point', 'end_point') NOT NULL);
CREATE TABLE text_values (value TEXT NOT NULL);
