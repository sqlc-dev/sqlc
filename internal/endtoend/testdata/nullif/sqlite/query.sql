CREATE TABLE IF NOT EXISTS Child(
    ID INTEGER NOT NULL PRIMARY KEY,
    PublicID TEXT NOT NULL,
    Name TEXT NOT NULL,
    DOB TEXT NOT NULL, 
    Photo BLOB
);

-- name: UpdateChildNullIf :one
UPDATE Child SET
    name  = COALESCE(nullif(?,''), name),
    dob   = COALESCE(nullif(?,''), dob),
    photo = COALESCE(?,photo)
    WHERE
        publicid = ?
    RETURNING *;

-- name: UpdateChild :one
UPDATE Child SET
    name  = COALESCE(?, name),
    dob   = COALESCE(?, dob),
    photo = COALESCE(?, photo)
    WHERE
        publicid = ?
    RETURNING *;