/* name: FooLimit :many */
SELECT a FROM foo
LIMIT ?;

/* name: FooLimitOffset :many */
SELECT a FROM foo
LIMIT ? OFFSET ?;
