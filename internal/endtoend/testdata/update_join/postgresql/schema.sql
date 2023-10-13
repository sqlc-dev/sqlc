CREATE TABLE primary_table (
  id         INT PRIMARY KEY,
  user_id    INT NOT NULL
);

CREATE TABLE join_table (
  id                INT PRIMARY KEY,
  primary_table_id  INT NOT NULL,
  other_table_id    INT NOT NULL,
  is_active         BOOLEAN NOT NULL
);

