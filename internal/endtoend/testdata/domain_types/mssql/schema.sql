CREATE TYPE dbo.PhoneNumber FROM varchar(20) NOT NULL;
CREATE TYPE Code FROM char(3);
CREATE SCHEMA sales;
CREATE TYPE sales.Amount FROM decimal(10, 2);

CREATE TABLE contacts (
  id BIGINT IDENTITY(1,1) PRIMARY KEY,
  phone dbo.PhoneNumber,
  code Code,
  total sales.Amount
);
