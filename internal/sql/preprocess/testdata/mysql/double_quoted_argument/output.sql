SELECT
    (CASE WHEN ? = "calendar" THEN cal.Timezone ELSE ? END),
    COALESCE(?, "")
FROM Calendar cal;
