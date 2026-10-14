CREATE TABLE venues (
  id INT64 NOT NULL,
) PRIMARY KEY (id);
CREATE TABLE arenas (
  id INT64 NOT NULL,
) PRIMARY KEY (id);
ALTER TABLE venues RENAME TO venue;
RENAME TABLE arenas TO stadiums, stadiums TO arena;
