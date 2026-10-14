-- name: AlbumIDs :one
SELECT singer_id, ARRAY(SELECT album_id FROM albums WHERE albums.singer_id = singers.singer_id) AS album_ids
FROM singers WHERE singer_id = @singer_id;

-- name: Literals :one
SELECT [1, 2, 3] AS ints, ARRAY<STRING>['a'] AS strs, DATE '2024-01-02' AS d,
  TIMESTAMP '2024-01-02 03:04:05+00' AS ts, NUMERIC '1.5' AS n, JSON '{"a": 1}' AS j,
  b'abc' AS raw, INTERVAL 3 DAY AS iv;
