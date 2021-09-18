CREATE TABLE foo (a text, b text);
CREATE TABLE bar (a text, b text);
/* name: FooByBarB :many */
SELECT a, b from foo where foo.a in (select a from bar where bar.b = ?);

/* name: FooByList :many */
SELECT a, b from foo where foo.a in (?, ?);

/* name: FooByNotList :many */
SELECT a, b from foo where foo.a not in (?, ?);

/* name: FooByParamList :many */
SELECT a, b from foo where ? in (foo.a, foo.b);
