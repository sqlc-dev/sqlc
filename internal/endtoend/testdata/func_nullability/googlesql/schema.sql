CREATE TABLE events (
  id INT64 NOT NULL,
  name STRING(MAX) NOT NULL,
  nickname STRING(MAX),
  amount NUMERIC NOT NULL,
  ratio FLOAT64,
  happened_at TIMESTAMP NOT NULL,
  details JSON NOT NULL,
  scores ARRAY<INT64> NOT NULL,
) PRIMARY KEY (id);
