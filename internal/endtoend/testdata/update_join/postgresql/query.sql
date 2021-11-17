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

-- name: UpdateJoin :exec
UPDATE  join_table
SET     is_active = $1
FROM    primary_table
WHERE   join_table.id = $2
        AND primary_table.user_id = $3
        AND join_table.primary_table_id = primary_table.id;
