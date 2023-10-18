-- name: TestFuncGetTime :one
select * from test_func_get_time();

-- name: TestFuncSelectString :one
select * from test_func_select_string($1);

-- name: TestFuncSelectBlog :one
select * from test_func_select_blog($1);
