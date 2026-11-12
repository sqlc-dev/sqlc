CREATE TABLE sales (
    id INT NOT NULL PRIMARY KEY,
    author_id INT NOT NULL,
    year INT NOT NULL,
    amount INT NOT NULL,
    note NVARCHAR(50) NULL
);

CREATE TABLE stats (
    id INT NOT NULL PRIMARY KEY,
    name VARCHAR(20) NULL,
    q1 INT NOT NULL,
    q2 INT NULL,
    q3 INT NULL
);
