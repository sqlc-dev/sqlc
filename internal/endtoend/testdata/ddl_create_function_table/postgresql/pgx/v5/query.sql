-- name: SelectAll :many
SELECT * FROM public.get_test();

-- name: SelectWithTime :many
SELECT * FROM public.get_test($1::timestamp);


-- name: GetTestIDByMessageFields :one
SELECT test_id
FROM public.get_all_tests_at_moment (p_target_time => $1)
WHERE test_number = $2
  AND departure_test_code = $3
  AND destination_test_code = $4
  AND index_number = $5;

