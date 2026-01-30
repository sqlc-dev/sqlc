
-- name: AddAuthor :execlastid
INSERT INTO authors (
    address,
    name,
    bio
) VALUES (
    $address,
    COALESCE($cal_name, ''),
    COALESCE($cal_description, '')
);
