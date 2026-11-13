CREATE TABLE people (
    id INT IDENTITY(1,1) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    nickname NVARCHAR(20),
    code CHAR(3) NOT NULL,
    score DECIMAL(10,2) NOT NULL,
    joined DATETIME2(3) NOT NULL
);

CREATE TABLE orgs (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    title NVARCHAR(100) NOT NULL,
    code VARCHAR(5),
    budget MONEY NOT NULL,
    founded DATE NOT NULL
);

CREATE TABLE ledger (
    qty INT NOT NULL,
    flag BIT NOT NULL,
    amount DECIMAL(5,2) NOT NULL
);
