CREATE TABLE sales (
  id INT IDENTITY(1,1) PRIMARY KEY,
  royalties DECIMAL(10,2),
  units INT,
  refunds SMALLINT,
  flags TINYINT,
  total BIGINT,
  ratio REAL,
  weight FLOAT,
  fee MONEY,
  small_fee SMALLMONEY,
  code CHAR(4),
  title VARCHAR(50),
  ncode NCHAR(4),
  name NVARCHAR(50)
);
