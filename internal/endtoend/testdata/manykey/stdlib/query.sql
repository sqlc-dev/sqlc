CREATE TABLE foo (
	group_id INT NOT NULL,
	score INT NOT NULL
);

-- name: SelectScoreSums :many key=group_id
SELECT group_id, SUM(score) FROM foo GROUP BY group_id;
