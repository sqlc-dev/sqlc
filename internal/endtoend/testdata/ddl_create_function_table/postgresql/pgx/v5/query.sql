-- name: SelectAll :many
SELECT * FROM public.get_test();

-- name: SelectWithTime :many
SELECT * FROM public.get_test($1::timestamp);
