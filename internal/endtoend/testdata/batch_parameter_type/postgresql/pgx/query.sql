-- name: InsertMappping :batchexec
WITH
    table1
        AS (
            SELECT
                version
            FROM
                solar_commcard_mapping
            WHERE
                "deviceId" = $1
            ORDER BY
                "updatedAt" DESC
            LIMIT
                1
        )
INSERT
INTO
    solar_commcard_mapping
        ("deviceId", version, sn, "updatedAt")
SELECT
    $1, @version::text, $3, $4
WHERE
    NOT
        EXISTS(
            SELECT
                *
            FROM
                table1
            WHERE
                table1.version = @version::text
        )
    OR NOT EXISTS(SELECT * FROM table1);
