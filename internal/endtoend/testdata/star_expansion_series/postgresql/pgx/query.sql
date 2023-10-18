-- name: CountAlertReportBy :many
select DATE_TRUNC($1,ts)::text as datetime,coalesce(count,0) as count from 
(
    SELECT  DATE_TRUNC($1,eventdate) as hr ,count(*)
    FROM    alertreport
    where eventdate between $2 and $3
    GROUP BY   1
) AS cnt 
right outer join ( SELECT * FROM generate_series ( $2, $3, CONCAT('1 ',$1)::interval) AS ts ) as dte
on DATE_TRUNC($1, ts ) = cnt.hr 
order by 1 asc;
