-- name: SelectAllColsFt :many
SELECT b FROM ft
WHERE b MATCH ?;

-- name: SelectAllColsTblFt :many
SELECT b, c FROM tbl_ft
WHERE b MATCH ?;

-- name: SelectOneColFt :many
SELECT b FROM ft
WHERE b = ?;

-- name: SelectOneColTblFt :many
SELECT c FROM tbl_ft
WHERE b = ?;

-- name: SelectHightlighFunc :many
SELECT highlight(tbl_ft, 0, '<b>', '</b>') FROM tbl_ft
WHERE b MATCH ?;

-- name: SelectSnippetFunc :many
SELECT snippet(tbl_ft, 0, '<b>', '</b>', 'aa', ?) FROM tbl_ft;

-- name: SelectBm25Func :many
SELECT *, bm25(tbl_ft, 2.0) FROM tbl_ft
WHERE b MATCH ? ORDER BY bm25(tbl_ft);

-- name: UpdateTblFt :exec
UPDATE tbl_ft SET c = ? WHERE b = ?;

-- name: DeleteTblFt :exec
DELETE FROM tbl_ft WHERE b = ?;

-- name: InsertTblFt :exec
INSERT INTO tbl_ft(b, c) VALUES(?, ?);
