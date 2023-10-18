-- name: TestFuncSelectBlog :many
select id, name from test_select_blog($1);
