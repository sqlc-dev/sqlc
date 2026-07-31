SELECT
    (CASE WHEN sqlc.arg("Timezone") = "calendar" THEN cal.Timezone ELSE sqlc.arg("Timezone") END),
    COALESCE(sqlc.narg("calName"), "")
FROM Calendar cal;
