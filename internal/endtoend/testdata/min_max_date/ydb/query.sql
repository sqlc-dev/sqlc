-- name: ActivityStats :one
SELECT COUNT(*) AS NumOfActivities,
        CAST(MIN(event_time) AS Timestamp) AS MinDate, 
        CAST(MAX(event_time) AS Timestamp) AS MaxDate 
FROM activities 
WHERE account_id = $account_id;

