-- name: Grouping :many
SELECT g, GROUPING(g) AS gr, GROUPING_ID(g, id) AS gid, APPROX_COUNT_DISTINCT(id) AS approx, CHECKSUM_AGG(id) AS ca
FROM s GROUP BY ROLLUP(g, id);

-- name: Percentiles :many
SELECT PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY v) OVER (PARTITION BY g) AS med FROM s;

-- name: Metadata :one
SELECT DB_NAME() AS d, SCHEMA_NAME() AS sn, OBJECT_ID(N'dbo.s') AS o, OBJECT_NAME(@oid) AS obj,
    IDENT_CURRENT(N's') AS ic, ROWCOUNT_BIG() AS rc, ORIGINAL_LOGIN() AS ol, DB_ID() AS dbid,
    ERROR_MESSAGE() AS em, ERROR_NUMBER() AS en, XACT_STATE() AS xs, SESSION_CONTEXT(N'k') AS sc;

-- name: Strings :many
SELECT STRING_ESCAPE(g, 'json') AS se, SOUNDEX(g) AS sx, DIFFERENCE(g, @other) AS df,
    LTRIM(g, N'x') AS lt, RTRIM(g, N'x') AS rt, COMPRESS(g) AS c, REGEXP_REPLACE(g, @pattern, N'b') AS rr,
    REGEXP_SUBSTR(g, N'a') AS rs, REGEXP_COUNT(g, N'a') AS rc
FROM s;

-- name: Bits :many
SELECT BIT_COUNT(id) AS bc, LEFT_SHIFT(id, 1) AS ls, GET_BIT(id, 1) AS gb, SET_BIT(id, @bit) AS sb,
    DATE_BUCKET(week, 2, CAST(SYSDATETIME() AS date)) AS db
FROM s;
