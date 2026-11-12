CREATE TABLE books (
  id INT IDENTITY(1,1) PRIMARY KEY,
  title VARCHAR(100) NOT NULL,
  subtitle VARCHAR(100),
  name NVARCHAR(100) NOT NULL,
  note NVARCHAR(200),
  price DECIMAL(10,2) NOT NULL,
  discount DECIMAL(5,2),
  qty INT NOT NULL,
  rating INT,
  weight FLOAT,
  cost MONEY,
  published DATE,
  updated DATETIME2(3),
  offset_at DATETIMEOFFSET,
  doc NVARCHAR(MAX)
);
