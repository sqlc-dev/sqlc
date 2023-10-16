-- name: ActivityStats :one
SELECT COUNT(*) as NumOfActivities,
        MIN(event_time) as MinDate, 
        MAX(event_time) as MaxDate 
FROM activities 
WHERE account_id = $1;
