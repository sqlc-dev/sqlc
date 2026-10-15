CREATE TABLE singers (
  singer_id INT64 NOT NULL,
  name STRING(MAX) NOT NULL,
  tags ARRAY<STRING(MAX)>,
  scores ARRAY<INT64> NOT NULL,
) PRIMARY KEY (singer_id);

CREATE TABLE albums (
  singer_id INT64 NOT NULL,
  album_id INT64 NOT NULL,
  title STRING(MAX),
) PRIMARY KEY (singer_id, album_id),
  INTERLEAVE IN PARENT singers ON DELETE CASCADE;
