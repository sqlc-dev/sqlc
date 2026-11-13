-- name: Literals :many
SELECT N'hello' AS a1, N'' AS a2, 'hi' AS a3, a + N'xy' AS a4, N'😀' AS a5 FROM p WHERE a LIKE N'%x%' AND b = @b;
