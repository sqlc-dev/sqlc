-- name: GroupConcat :many
SELECT student_name, GROUP_CONCAT(test_score)
FROM student
GROUP BY student_name;

-- name: GroupConcatOrderBy :many
SELECT student_name,
    GROUP_CONCAT(DISTINCT test_score ORDER BY test_score DESC SEPARATOR ' ')
FROM student
GROUP BY student_name;
