-- name: AddEvent :execlastid
INSERT INTO `Event` (
    Timezone
) VALUES (
    (CASE WHEN sqlc.arg("Timezone") = "calendar" THEN (SELECT cal.Timezone FROM Calendar cal WHERE cal.IdKey = sqlc.arg("calendarIdKey")) ELSE sqlc.arg("Timezone") END)
);

-- name: AddAuthor :execlastid
INSERT INTO authors (
    address,
    name,
    bio
) VALUES (
    ?,
    COALESCE(sqlc.narg("calName"), ""),
    COALESCE(sqlc.narg("calDescription"), "")
);
