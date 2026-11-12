CREATE TABLE readings (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    raw VARCHAR(10) NOT NULL,
    note VARCHAR(10),
    amount DECIMAL(10,2)
);
