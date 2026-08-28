### case_named_params/mysql
reason: semantic drift in query.sql
--- query diff:
@@ -1,6 +1,6 @@
 -- name: ListAuthors :one
-SELECT  *
-FROM    authors
-WHERE   email = CASE WHEN sqlc.arg(email) = '' then NULL else sqlc.arg(email) END
-        OR username = CASE WHEN sqlc.arg(username) = '' then NULL else sqlc.arg(username) END 
-LIMIT   1;
+SELECT *
+FROM authors
+WHERE email = CASE WHEN sqlc.arg(email) = '' THEN '' ELSE sqlc.arg(email) END
+  OR username = CASE WHEN sqlc.arg(username) = '' THEN '' ELSE sqlc.arg(username) END
+LIMIT 1;
--- semantic check (query.sql):
stmt 1 differs:
  orig: SELECT * FROM `authors` WHERE `email`=CASE WHEN `sqlc`.`arg`(`email`)=_UTF8MB4'' THEN NULL ELSE `sqlc`.`arg`(`email`) END OR `username`=CASE WHEN `sqlc`.`arg`(`username`)=_UTF8MB4'' THEN NULL ELSE `sqlc`.`arg`(`username`) END LIMIT 1
  fmt:  SELECT * FROM `authors` WHERE `email`=CASE WHEN `sqlc`.`arg`(`email`)=_UTF8MB4'' THEN _UTF8MB4'' ELSE `sqlc`.`arg`(`email`) END OR `username`=CASE WHEN `sqlc`.`arg`(`username`)=_UTF8MB4'' THEN _UTF8MB4'' ELSE `sqlc`.`arg`(`username`) END LIMIT 1


### case_value_param/mysql
reason: semantic drift in query.sql
--- query diff:
@@ -1,3 +1,3 @@
 -- name: Update :exec
 UPDATE testing
-SET value = CASE ? WHEN true THEN 'Hello' WHEN false THEN 'Goodbye' ELSE value END;
+SET value = CASE ? WHEN 1 THEN 'Hello' WHEN 0 THEN 'Goodbye' ELSE value END;
--- semantic check (query.sql):
stmt 1 differs:
  orig: UPDATE `testing` SET `value`=CASE ? WHEN TRUE THEN _UTF8MB4'Hello' WHEN FALSE THEN _UTF8MB4'Goodbye' ELSE `value` END
  fmt:  UPDATE `testing` SET `value`=CASE ? WHEN 1 THEN _UTF8MB4'Hello' WHEN 0 THEN _UTF8MB4'Goodbye' ELSE `value` END


### cte_count/mysql
reason: semantic drift in query.sql
--- query diff:
@@ -1,8 +1,8 @@
 -- name: CTECount :many
-WITH all_count AS (
-	SELECT count(*) FROM bar
-), ready_count AS (
-	SELECT count(*) FROM bar WHERE ready = true
+WITH all_count AS (SELECT count(*) FROM bar),
+ready_count AS (
+  SELECT count(*) FROM bar WHERE ready = 1
 )
 SELECT all_count.count, ready_count.count
-FROM all_count, ready_count;
+FROM all_count
+JOIN ready_count;
--- semantic check (query.sql):
stmt 1 differs:
  orig: WITH `all_count` AS (SELECT COUNT(1) FROM `bar`), `ready_count` AS (SELECT COUNT(1) FROM `bar` WHERE `ready`=TRUE) SELECT `all_count`.`count`,`ready_count`.`count` FROM (`all_count`) JOIN `ready_count`
  fmt:  WITH `all_count` AS (SELECT COUNT(1) FROM `bar`), `ready_count` AS (SELECT COUNT(1) FROM `bar` WHERE `ready`=1) SELECT `all_count`.`count`,`ready_count`.`count` FROM `all_count` JOIN `ready_count`


### cte_in_delete/mysql
reason: semantic drift in query.sql
--- query diff:
@@ -1,5 +1,4 @@
 -- name: DeleteReadyWithCTE :exec
-WITH ready_ids AS (
-	SELECT id FROM bar WHERE ready
-)
-DELETE FROM bar WHERE id IN (SELECT * FROM ready_ids);
+WITH ready_ids AS (SELECT id FROM bar WHERE ready)
+DELETE FROM bar
+WHERE id IN ((SELECT * FROM ready_ids));
--- semantic check (query.sql):
stmt 1 differs:
  orig: WITH `ready_ids` AS (SELECT `id` FROM `bar` WHERE `ready`) DELETE FROM `bar` WHERE `id` IN (SELECT * FROM `ready_ids`)
  fmt:  WITH `ready_ids` AS (SELECT `id` FROM `bar` WHERE `ready`) DELETE FROM `bar` WHERE `id` IN ((SELECT * FROM `ready_ids`))


### cte_recursive/mysql
reason: semantic drift in query.sql
--- query diff:
@@ -1,9 +1,12 @@
 -- name: CTERecursive :many
 WITH RECURSIVE cte AS (
-        SELECT b.* FROM bar AS b
-        WHERE b.id = ?
-    UNION ALL
-        SELECT b.*
-        FROM bar AS b, cte AS c
-        WHERE b.parent_id = c.id
-) SELECT * FROM cte;
+  SELECT b.*
+  FROM bar AS b
+  WHERE b.id = ?
+  UNION ALL
+  SELECT b.*
+  FROM bar AS b
+  JOIN cte AS c
+  WHERE b.parent_id = c.id
+)
+SELECT * FROM cte;
--- semantic check (query.sql):
stmt 1 differs:
  orig: WITH RECURSIVE `cte` AS (SELECT `b`.* FROM `bar` AS `b` WHERE `b`.`id`=? UNION ALL SELECT `b`.* FROM (`bar` AS `b`) JOIN `cte` AS `c` WHERE `b`.`parent_id`=`c`.`id`) SELECT * FROM `cte`
  fmt:  WITH RECURSIVE `cte` AS (SELECT `b`.* FROM `bar` AS `b` WHERE `b`.`id`=? UNION ALL SELECT `b`.* FROM `bar` AS `b` JOIN `cte` AS `c` WHERE `b`.`parent_id`=`c`.`id`) SELECT * FROM `cte`


### exec_lastid/go_postgresql_stdlib
reason: semantic drift in query.sql
--- query diff:
@@ -1,2 +1,2 @@
 -- name: InsertBar :execlastid
-INSERT INTO bar () VALUES ();
\ No newline at end of file
+INSERT INTO bar VALUES ();
--- semantic check (query.sql):
stmt 1 differs:
  orig: INSERT INTO `bar` () VALUES ()
  fmt:  INSERT INTO `bar` VALUES ()


### identifier_case_sensitivity
reason: semantic drift in query.sql
--- query diff:
@@ -1,18 +1,17 @@
 -- name: GetAuthor :one
-SELECT * FROM Authors
-WHERE ID = ? LIMIT 1;
+SELECT *
+FROM authors
+WHERE id = ?
+LIMIT 1;
 
 -- name: ListAuthors :many
-SELECT * FROM Authors
-ORDER BY Name;
+SELECT *
+FROM authors
+ORDER BY name;
 
 -- name: CreateAuthor :execresult
-INSERT INTO Authors (
-  Name, Bio
-) VALUES (
-  ?, ?
-);
+INSERT INTO authors (name, bio)
+VALUES (?, ?);
 
 -- name: DeleteAuthor :exec
-DELETE FROM Authors
-WHERE ID = ?;
+DELETE FROM authors WHERE id = ?;
--- semantic check (query.sql):
stmt 1 differs:
  orig: SELECT * FROM `Authors` WHERE `ID`=? LIMIT 1
  fmt:  SELECT * FROM `authors` WHERE `id`=? LIMIT 1
stmt 2 differs:
  orig: SELECT * FROM `Authors` ORDER BY `Name`
  fmt:  SELECT * FROM `authors` ORDER BY `name`
stmt 3 differs:
  orig: INSERT INTO `Authors` (`Name`,`Bio`) VALUES (?,?)
  fmt:  INSERT INTO `authors` (`name`,`bio`) VALUES (?,?)
stmt 4 differs:
  orig: DELETE FROM `Authors` WHERE `ID`=?
  fmt:  DELETE FROM `authors` WHERE `id`=?


### in_union/mysql
reason: semantic drift in query.sql
--- query diff:
@@ -1,3 +1,4 @@
 -- name: GetAuthors :many
-SELECT * FROM authors
-WHERE author_id IN (SELECT author_id FROM book1 UNION SELECT author_id FROM book2);
+SELECT *
+FROM authors
+WHERE author_id IN ((SELECT author_id FROM book1 UNION SELECT author_id FROM book2));
--- semantic check (query.sql):
stmt 1 differs:
  orig: SELECT * FROM `authors` WHERE `author_id` IN (SELECT `author_id` FROM `book1` UNION SELECT `author_id` FROM `book2`)
  fmt:  SELECT * FROM `authors` WHERE `author_id` IN ((SELECT `author_id` FROM `book1` UNION SELECT `author_id` FROM `book2`))


### join_from/mysql
reason: semantic drift in query.sql
--- query diff:
@@ -1,2 +1,2 @@
 -- name: MultiFrom :many
-SELECT email FROM bar, foo WHERE login = ?;
+SELECT email FROM bar JOIN foo WHERE login = ?;
--- semantic check (query.sql):
stmt 1 differs:
  orig: SELECT `email` FROM (`bar`) JOIN `foo` WHERE `login`=?
  fmt:  SELECT `email` FROM `bar` JOIN `foo` WHERE `login`=?


### join_left/mysql
reason: semantic drift in query.sql
--- query diff:
@@ -1,59 +1,52 @@
 -- name: GetMayors :many
 SELECT
-    user_id,
-    mayors.full_name
+  user_id,
+  mayors.full_name
 FROM users
 LEFT JOIN cities USING (city_id)
-INNER JOIN mayors USING (mayor_id);
+JOIN mayors USING (mayor_id);
 
 -- name: GetMayorsOptional :many
 SELECT
-    user_id,
-    cities.city_id,
-    mayors.full_name
+  user_id,
+  cities.city_id,
+  mayors.full_name
 FROM users
 LEFT JOIN cities USING (city_id)
 LEFT JOIN mayors USING (mayor_id);
 
 -- name: AllAuthors :many
-SELECT  *
-FROM    authors a
-        LEFT JOIN authors p
-            ON a.parent_id = p.id;
+SELECT *
+FROM authors AS a
+LEFT JOIN authors AS p ON a.parent_id = p.id;
 
 -- name: AllAuthorsAliases :many
-SELECT  *
-FROM    authors a
-        LEFT JOIN authors p
-            ON a.parent_id = p.id;
+SELECT *
+FROM authors AS a
+LEFT JOIN authors AS p ON a.parent_id = p.id;
 
 -- name: AllAuthorsAliases2 :many
-SELECT  a.*, p.*
-FROM    authors a
-        LEFT JOIN authors p
-            ON a.parent_id = p.id;
+SELECT a.*, p.*
+FROM authors AS a
+LEFT JOIN authors AS p ON a.parent_id = p.id;
 
 -- name: AllSuperAuthors :many
-SELECT  *
-FROM    authors
-        LEFT JOIN super_authors
-            ON authors.parent_id = super_authors.super_id;
+SELECT *
+FROM authors
+LEFT JOIN super_authors ON authors.parent_id = super_authors.super_id;
 
 -- name: AllSuperAuthorsAliases :many
-SELECT  *
-FROM    authors a
-        LEFT JOIN super_authors sa
-            ON a.parent_id = sa.super_id;
+SELECT *
+FROM authors AS a
+LEFT JOIN super_authors AS sa ON a.parent_id = sa.super_id;
 
 -- name: AllSuperAuthorsAliases2 :many
-SELECT  a.*, sa.*
-FROM    authors a
-        LEFT JOIN super_authors sa
-            ON a.parent_id = sa.super_id;
+SELECT a.*, sa.*
+FROM authors AS a
+LEFT JOIN super_authors AS sa ON a.parent_id = sa.super_id;
 
 -- name: GetSuggestedUsersByID :many
-SELECT  DISTINCT u.*, m.*
-FROM    users_2 u
-        LEFT JOIN media m
-            ON u.user_avatar_id = m.media_id
-WHERE   u.user_id != @user_id;
+SELECT u.*, m.*
+FROM users_2 AS u
+LEFT JOIN media AS m ON u.user_avatar_id = m.media_id
+WHERE u.user_id != @user_id;
--- semantic check (query.sql):
stmt 9 differs:
  orig: SELECT DISTINCT `u`.*,`m`.* FROM `users_2` AS `u` LEFT JOIN `media` AS `m` ON `u`.`user_avatar_id`=`m`.`media_id` WHERE `u`.`user_id`!=@`user_id`
  fmt:  SELECT `u`.*,`m`.* FROM `users_2` AS `u` LEFT JOIN `media` AS `m` ON `u`.`user_avatar_id`=`m`.`media_id` WHERE `u`.`user_id`!=@`user_id`


### mysql_optimizer_hints/mysql
reason: semantic drift in query.sql
--- query diff:
@@ -1,9 +1,7 @@
 -- name: InlineHint :one
-SELECT /*+ MAX_EXECUTION_TIME(1000) */ bar FROM foo LIMIT 1;
+SELECT bar FROM foo LIMIT 1;
 
 -- name: MultilineHint :one
-SELECT
-/*+ MAX_EXECUTION_TIME(1000) */
-bar
+SELECT bar
 FROM foo
 LIMIT 1;
--- semantic check (query.sql):
stmt 1 differs:
  orig: SELECT /*+ MAX_EXECUTION_TIME(1000)*/ `bar` FROM `foo` LIMIT 1
  fmt:  SELECT `bar` FROM `foo` LIMIT 1
stmt 2 differs:
  orig: SELECT /*+ MAX_EXECUTION_TIME(1000)*/ `bar` FROM `foo` LIMIT 1
  fmt:  SELECT `bar` FROM `foo` LIMIT 1


### order_by_union/mysql
reason: semantic drift in query.sql
--- query diff:
@@ -1,5 +1,4 @@
 -- name: ListAuthorsUnion :many
-SELECT name as foo FROM authors
+SELECT name AS foo FROM authors
 UNION
-SELECT first_name as foo FROM people
-ORDER BY foo;
+SELECT first_name AS foo FROM people;
--- semantic check (query.sql):
stmt 1 differs:
  orig: SELECT `name` AS `foo` FROM `authors` UNION SELECT `first_name` AS `foo` FROM `people` ORDER BY `foo`
  fmt:  SELECT `name` AS `foo` FROM `authors` UNION SELECT `first_name` AS `foo` FROM `people`


### overrides_go_types/mysql
reason: semantic drift in query.sql
--- query diff:
@@ -1,2 +1,2 @@
 -- name: TestIN :many
-SELECT * FROM foo WHERE retyped IN (sqlc.slice(paramName));
+SELECT * FROM foo WHERE retyped IN (sqlc.slice(paramname));
--- semantic check (query.sql):
stmt 1 differs:
  orig: SELECT * FROM `foo` WHERE `retyped` IN (`sqlc`.`slice`(`paramName`))
  fmt:  SELECT * FROM `foo` WHERE `retyped` IN (`sqlc`.`slice`(`paramname`))


### params_in_nested_func/mysql
reason: semantic drift in query.sql
--- query diff:
@@ -1,9 +1,7 @@
 -- name: GetGroups :many
 SELECT
-    rg.groupId,
-    rg.groupName
-FROM 
-    RouterGroup rg
-WHERE
-    rg.groupName LIKE CONCAT('%', COALESCE(sqlc.narg('groupName'), rg.groupName), '%') AND
-    rg.groupId = COALESCE(sqlc.narg('groupId'), rg.groupId);
+  rg.groupid,
+  rg.groupname
+FROM routergroup AS rg
+WHERE rg.groupname LIKE concat('%', coalesce(sqlc.narg('groupName'), rg.groupname), '%')
+  AND rg.groupid = coalesce(sqlc.narg('groupId'), rg.groupid);
--- semantic check (query.sql):
stmt 1 differs:
  orig: SELECT `rg`.`groupId`,`rg`.`groupName` FROM `RouterGroup` AS `rg` WHERE `rg`.`groupName` LIKE CONCAT(_UTF8MB4'%', COALESCE(`sqlc`.`narg`(_UTF8MB4'groupName'), `rg`.`groupName`), _UTF8MB4'%') AND `rg`.`groupId`=COALESCE(`sqlc`.`narg`(_UTF8MB4'groupId'), `rg`.`groupId`)
  fmt:  SELECT `rg`.`groupid`,`rg`.`groupname` FROM `routergroup` AS `rg` WHERE `rg`.`groupname` LIKE CONCAT(_UTF8MB4'%', COALESCE(`sqlc`.`narg`(_UTF8MB4'groupName'), `rg`.`groupname`), _UTF8MB4'%') AND `rg`.`groupid`=COALESCE(`sqlc`.`narg`(_UTF8MB4'groupId'), `rg`.`groupid`)


### pattern_in_expr/mysql
reason: semantic drift in query.sql
--- query diff:
@@ -1,11 +1,11 @@
 /* name: FooByBarB :many */
-SELECT a, b from foo where foo.a in (select a from bar where bar.b = ?);
+SELECT a, b FROM foo WHERE foo.a IN ((SELECT a FROM bar WHERE bar.b = ?));
 
 /* name: FooByList :many */
-SELECT a, b from foo where foo.a in (?, ?);
+SELECT a, b FROM foo WHERE foo.a IN (?, ?);
 
 /* name: FooByNotList :many */
-SELECT a, b from foo where foo.a not in (?, ?);
+SELECT a, b FROM foo WHERE foo.a NOT IN (?, ?);
 
 /* name: FooByParamList :many */
-SELECT a, b from foo where ? in (foo.a, foo.b);
+SELECT a, b FROM foo WHERE ? IN (foo.a, foo.b);
--- semantic check (query.sql):
stmt 1 differs:
  orig: SELECT `a`,`b` FROM `foo` WHERE `foo`.`a` IN (SELECT `a` FROM `bar` WHERE `bar`.`b`=?)
  fmt:  SELECT `a`,`b` FROM `foo` WHERE `foo`.`a` IN ((SELECT `a` FROM `bar` WHERE `bar`.`b`=?))


### select_union/mysql
reason: semantic drift in query.sql
--- query diff:
@@ -7,7 +7,8 @@
 SELECT * FROM foo
 UNION
 SELECT * FROM foo
-LIMIT ? OFFSET ?;
+LIMIT ?
+OFFSET ?;
 
 -- name: SelectExcept :many
 SELECT * FROM foo
@@ -25,6 +26,6 @@
 SELECT * FROM bar;
 
 -- name: SelectUnionAliased :many
-(SELECT * FROM foo)
+SELECT * FROM foo
 UNION
 SELECT * FROM bar;
--- semantic check (query.sql):
stmt 6 differs:
  orig: (SELECT * FROM `foo`) UNION SELECT * FROM `bar`
  fmt:  SELECT * FROM `foo` UNION SELECT * FROM `bar`


### show_warnings/mysql
reason: semantic drift in query.sql
--- query diff:
@@ -1,2 +1,2 @@
 -- name: ShowWarnings :many
-SHOW WARNINGS;
+SELECT '' AS level, 0 AS code, '' AS message;
--- semantic check (query.sql):
stmt 1 differs:
  orig: SHOW WARNINGS
  fmt:  SELECT _UTF8MB4'' AS `level`,0 AS `code`,_UTF8MB4'' AS `message`


### sqlc_arg/mysql
reason: semantic drift in query.sql
--- query diff:
@@ -5,5 +5,5 @@
 SELECT name FROM foo WHERE name = sqlc.arg('slug');
 
 /* name: Complicated :many */
-WITH names AS (SELECT name from foo WHERE foo.name = sqlc.arg('slug'))
-SELECT name FROM names WHERE name IN (SELECT name FROM foo WHERE foo.name = sqlc.arg('slug'));
+WITH names AS (SELECT name FROM foo WHERE foo.name = sqlc.arg('slug'))
+SELECT name FROM names WHERE name IN ((SELECT name FROM foo WHERE foo.name = sqlc.arg('slug')));
--- semantic check (query.sql):
stmt 3 differs:
  orig: WITH `names` AS (SELECT `name` FROM `foo` WHERE `foo`.`name`=`sqlc`.`arg`(_UTF8MB4'slug')) SELECT `name` FROM `names` WHERE `name` IN (SELECT `name` FROM `foo` WHERE `foo`.`name`=`sqlc`.`arg`(_UTF8MB4'slug'))
  fmt:  WITH `names` AS (SELECT `name` FROM `foo` WHERE `foo`.`name`=`sqlc`.`arg`(_UTF8MB4'slug')) SELECT `name` FROM `names` WHERE `name` IN ((SELECT `name` FROM `foo` WHERE `foo`.`name`=`sqlc`.`arg`(_UTF8MB4'slug')))


### star_expansion_join/mysql
reason: semantic drift in query.sql
--- query diff:
@@ -1,2 +1,2 @@
 /* name: StarExpansionJoin :many */
-SELECT * FROM foo, bar;
+SELECT * FROM foo JOIN bar;
--- semantic check (query.sql):
stmt 1 differs:
  orig: SELECT * FROM (`foo`) JOIN `bar`
  fmt:  SELECT * FROM `foo` JOIN `bar`


### update_inner_join
reason: semantic drift in query.sql
--- query diff:
@@ -1,2 +1,2 @@
 -- name: UpdateXWithY :exec
-UPDATE x INNER JOIN y ON y.a = x.a SET x.b = y.b;
+UPDATE x, y SET b = y.b;
--- semantic check (query.sql):
stmt 1 differs:
  orig: UPDATE `x` JOIN `y` ON `y`.`a`=`x`.`a` SET `x`.`b`=`y`.`b`
  fmt:  UPDATE (`x`) JOIN `y` SET `b`=`y`.`b`


### update_join/mysql
reason: semantic drift in query.sql
--- query diff:
@@ -1,23 +1,17 @@
 -- name: UpdateJoin :exec
-UPDATE  join_table as jt
-        JOIN primary_table as pt
-            ON jt.primary_table_id = pt.id
-SET     jt.is_active = ?
-WHERE   jt.id = ?
-        AND pt.user_id = ?;
+UPDATE join_table AS jt, primary_table AS pt
+SET is_active = ?
+WHERE jt.id = ?
+  AND pt.user_id = ?;
 
 -- name: UpdateLeftJoin :exec
-UPDATE  join_table as jt
-        LEFT JOIN primary_table as pt
-            ON jt.primary_table_id = pt.id
-SET     jt.is_active = ?
-WHERE   jt.id = ?
-        AND pt.user_id = ?;
+UPDATE join_table AS jt, primary_table AS pt
+SET is_active = ?
+WHERE jt.id = ?
+  AND pt.user_id = ?;
 
 -- name: UpdateRightJoin :exec
-UPDATE  join_table as jt
-        RIGHT JOIN primary_table as pt
-            ON jt.primary_table_id = pt.id
-SET     jt.is_active = ?
-WHERE   jt.id = ?
-        AND pt.user_id = ?;
+UPDATE join_table AS jt, primary_table AS pt
+SET is_active = ?
+WHERE jt.id = ?
+  AND pt.user_id = ?;
--- semantic check (query.sql):
stmt 1 differs:
  orig: UPDATE `join_table` AS `jt` JOIN `primary_table` AS `pt` ON `jt`.`primary_table_id`=`pt`.`id` SET `jt`.`is_active`=? WHERE `jt`.`id`=? AND `pt`.`user_id`=?
  fmt:  UPDATE (`join_table` AS `jt`) JOIN `primary_table` AS `pt` SET `is_active`=? WHERE `jt`.`id`=? AND `pt`.`user_id`=?
stmt 2 differs:
  orig: UPDATE `join_table` AS `jt` LEFT JOIN `primary_table` AS `pt` ON `jt`.`primary_table_id`=`pt`.`id` SET `jt`.`is_active`=? WHERE `jt`.`id`=? AND `pt`.`user_id`=?
  fmt:  UPDATE (`join_table` AS `jt`) JOIN `primary_table` AS `pt` SET `is_active`=? WHERE `jt`.`id`=? AND `pt`.`user_id`=?
stmt 3 differs:
  orig: UPDATE `join_table` AS `jt` RIGHT JOIN `primary_table` AS `pt` ON `jt`.`primary_table_id`=`pt`.`id` SET `jt`.`is_active`=? WHERE `jt`.`id`=? AND `pt`.`user_id`=?
  fmt:  UPDATE (`join_table` AS `jt`) JOIN `primary_table` AS `pt` SET `is_active`=? WHERE `jt`.`id`=? AND `pt`.`user_id`=?


### update_two_table/mysql
reason: semantic drift in query.sql
--- query diff:
@@ -1,10 +1,6 @@
 -- name: DeleteAuthor :exec
-UPDATE
-  authors,
-  books
-SET
-  authors.deleted_at = now(),
-  books.deleted_at = now()
-WHERE
-  books.is_amazing = 1
-  AND authors.name = sqlc.arg(name);
\ No newline at end of file
+UPDATE authors, books
+SET deleted_at = now(),
+  deleted_at = now()
+WHERE books.is_amazing = 1
+  AND authors.name = sqlc.arg(name);
--- semantic check (query.sql):
stmt 1 differs:
  orig: UPDATE (`authors`) JOIN `books` SET `authors`.`deleted_at`=NOW(), `books`.`deleted_at`=NOW() WHERE `books`.`is_amazing`=1 AND `authors`.`name`=`sqlc`.`arg`(`name`)
  fmt:  UPDATE (`authors`) JOIN `books` SET `deleted_at`=NOW(), `deleted_at`=NOW() WHERE `books`.`is_amazing`=1 AND `authors`.`name`=`sqlc`.`arg`(`name`)


### vet_explain/mysql
reason: semantic drift in query.sql
--- query diff:
@@ -1,146 +1,218 @@
 -- name: SelectById :one
-SELECT id FROM debug
-WHERE id = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE id = ?
+LIMIT 1;
 
 -- name: SelectByCsmallint :one
-SELECT id FROM debug
-WHERE Csmallint = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE csmallint = ?
+LIMIT 1;
 
 -- name: SelectByCint :one
-SELECT id FROM debug
-WHERE Cint = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE cint = ?
+LIMIT 1;
 
 -- name: SelectByCinteger :one
-SELECT id FROM debug
-WHERE Cinteger = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE cinteger = ?
+LIMIT 1;
 
 -- name: SelectByCdecimal :one
-SELECT id FROM debug
-WHERE Cdecimal = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE cdecimal = ?
+LIMIT 1;
 
 -- name: SelectByCnumeric :one
-SELECT id FROM debug
-WHERE Cnumeric = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE cnumeric = ?
+LIMIT 1;
 
 -- name: SelectByCfloat :one
-SELECT id FROM debug
-WHERE Cfloat = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE cfloat = ?
+LIMIT 1;
 
 -- name: SelectByCreal :one
-SELECT id FROM debug
-WHERE Creal = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE creal = ?
+LIMIT 1;
 
 -- name: SelectByCdoubleprecision :one
-SELECT id FROM debug
-WHERE Cdoubleprecision = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE cdoubleprecision = ?
+LIMIT 1;
 
 -- name: SelectByCdouble :one
-SELECT id FROM debug
-WHERE Cdouble = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE cdouble = ?
+LIMIT 1;
 
 -- name: SelectByCdec :one
-SELECT id FROM debug
-WHERE Cdec = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE cdec = ?
+LIMIT 1;
 
 -- name: SelectByCfixed :one
-SELECT id FROM debug
-WHERE Cfixed = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE cfixed = ?
+LIMIT 1;
 
 -- name: SelectByCtinyint :one
-SELECT id FROM debug
-WHERE Ctinyint = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE ctinyint = ?
+LIMIT 1;
 
 -- name: SelectByCbool :one
-SELECT id FROM debug
-WHERE Cbool = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE cbool = ?
+LIMIT 1;
 
 -- name: SelectByCmediumint :one
-SELECT id FROM debug
-WHERE Cmediumint = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE cmediumint = ?
+LIMIT 1;
 
 -- name: SelectByCbit :one
-SELECT id FROM debug
-WHERE Cbit = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE cbit = ?
+LIMIT 1;
 
 -- name: SelectByCdate :one
-SELECT id FROM debug
-WHERE Cdate = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE cdate = ?
+LIMIT 1;
 
 -- name: SelectByCdatetime :one
-SELECT id FROM debug
-WHERE Cdatetime = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE cdatetime = ?
+LIMIT 1;
 
 -- name: SelectByCtimestamp :one
-SELECT id FROM debug
-WHERE Ctimestamp = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE ctimestamp = ?
+LIMIT 1;
 
 -- name: SelectByCtime :one
-SELECT id FROM debug
-WHERE Ctime = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE ctime = ?
+LIMIT 1;
 
 -- name: SelectByCyear :one
-SELECT id FROM debug
-WHERE Cyear = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE cyear = ?
+LIMIT 1;
 
 -- name: SelectByCchar :one
-SELECT id FROM debug
-WHERE Cchar = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE cchar = ?
+LIMIT 1;
 
 -- name: SelectByCvarchar :one
-SELECT id FROM debug
-WHERE Cvarchar = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE cvarchar = ?
+LIMIT 1;
 
 -- name: SelectByCbinary :one
-SELECT id FROM debug
-WHERE Cbinary = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE cbinary = ?
+LIMIT 1;
 
 -- name: SelectByCvarbinary :one
-SELECT id FROM debug
-WHERE Cvarbinary = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE cvarbinary = ?
+LIMIT 1;
 
 -- name: SelectByCtinyblob :one
-SELECT id FROM debug
-WHERE Ctinyblob = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE ctinyblob = ?
+LIMIT 1;
 
 -- name: SelectByCblob :one
-SELECT id FROM debug
-WHERE Cblob = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE cblob = ?
+LIMIT 1;
 
 -- name: SelectByCmediumblob :one
-SELECT id FROM debug
-WHERE Cmediumblob = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE cmediumblob = ?
+LIMIT 1;
 
 -- name: SelectByClongblob :one
-SELECT id FROM debug
-WHERE Clongblob = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE clongblob = ?
+LIMIT 1;
 
 -- name: SelectByCtinytext :one
-SELECT id FROM debug
-WHERE Ctinytext = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE ctinytext = ?
+LIMIT 1;
 
 -- name: SelectByCtext :one
-SELECT id FROM debug
-WHERE Ctext = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE ctext = ?
+LIMIT 1;
 
 -- name: SelectByCmediumtext :one
-SELECT id FROM debug
-WHERE Cmediumtext = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE cmediumtext = ?
+LIMIT 1;
 
 -- name: SelectByClongtext :one
-SELECT id FROM debug
-WHERE Clongtext = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE clongtext = ?
+LIMIT 1;
 
 -- name: SelectByCenum :one
-SELECT id FROM debug
-WHERE Cenum = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE cenum = ?
+LIMIT 1;
 
 -- name: SelectByCset :one
-SELECT id FROM debug
-WHERE Cset = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE cset = ?
+LIMIT 1;
 
 -- name: SelectByCjson :one
-SELECT id FROM debug
-WHERE Cjson = ? LIMIT 1;
+SELECT id
+FROM debug
+WHERE cjson = ?
+LIMIT 1;
 
 --
 --
--- semantic check (query.sql):
stmt 2 differs:
  orig: SELECT `id` FROM `debug` WHERE `Csmallint`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `csmallint`=? LIMIT 1
stmt 3 differs:
  orig: SELECT `id` FROM `debug` WHERE `Cint`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `cint`=? LIMIT 1
stmt 4 differs:
  orig: SELECT `id` FROM `debug` WHERE `Cinteger`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `cinteger`=? LIMIT 1
stmt 5 differs:
  orig: SELECT `id` FROM `debug` WHERE `Cdecimal`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `cdecimal`=? LIMIT 1
stmt 6 differs:
  orig: SELECT `id` FROM `debug` WHERE `Cnumeric`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `cnumeric`=? LIMIT 1
stmt 7 differs:
  orig: SELECT `id` FROM `debug` WHERE `Cfloat`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `cfloat`=? LIMIT 1
stmt 8 differs:
  orig: SELECT `id` FROM `debug` WHERE `Creal`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `creal`=? LIMIT 1
stmt 9 differs:
  orig: SELECT `id` FROM `debug` WHERE `Cdoubleprecision`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `cdoubleprecision`=? LIMIT 1
stmt 10 differs:
  orig: SELECT `id` FROM `debug` WHERE `Cdouble`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `cdouble`=? LIMIT 1
stmt 11 differs:
  orig: SELECT `id` FROM `debug` WHERE `Cdec`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `cdec`=? LIMIT 1
stmt 12 differs:
  orig: SELECT `id` FROM `debug` WHERE `Cfixed`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `cfixed`=? LIMIT 1
stmt 13 differs:
  orig: SELECT `id` FROM `debug` WHERE `Ctinyint`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `ctinyint`=? LIMIT 1
stmt 14 differs:
  orig: SELECT `id` FROM `debug` WHERE `Cbool`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `cbool`=? LIMIT 1
stmt 15 differs:
  orig: SELECT `id` FROM `debug` WHERE `Cmediumint`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `cmediumint`=? LIMIT 1
stmt 16 differs:
  orig: SELECT `id` FROM `debug` WHERE `Cbit`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `cbit`=? LIMIT 1
stmt 17 differs:
  orig: SELECT `id` FROM `debug` WHERE `Cdate`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `cdate`=? LIMIT 1
stmt 18 differs:
  orig: SELECT `id` FROM `debug` WHERE `Cdatetime`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `cdatetime`=? LIMIT 1
stmt 19 differs:
  orig: SELECT `id` FROM `debug` WHERE `Ctimestamp`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `ctimestamp`=? LIMIT 1
stmt 20 differs:
  orig: SELECT `id` FROM `debug` WHERE `Ctime`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `ctime`=? LIMIT 1
stmt 21 differs:
  orig: SELECT `id` FROM `debug` WHERE `Cyear`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `cyear`=? LIMIT 1
stmt 22 differs:
  orig: SELECT `id` FROM `debug` WHERE `Cchar`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `cchar`=? LIMIT 1
stmt 23 differs:
  orig: SELECT `id` FROM `debug` WHERE `Cvarchar`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `cvarchar`=? LIMIT 1
stmt 24 differs:
  orig: SELECT `id` FROM `debug` WHERE `Cbinary`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `cbinary`=? LIMIT 1
stmt 25 differs:
  orig: SELECT `id` FROM `debug` WHERE `Cvarbinary`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `cvarbinary`=? LIMIT 1
stmt 26 differs:
  orig: SELECT `id` FROM `debug` WHERE `Ctinyblob`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `ctinyblob`=? LIMIT 1
stmt 27 differs:
  orig: SELECT `id` FROM `debug` WHERE `Cblob`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `cblob`=? LIMIT 1
stmt 28 differs:
  orig: SELECT `id` FROM `debug` WHERE `Cmediumblob`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `cmediumblob`=? LIMIT 1
stmt 29 differs:
  orig: SELECT `id` FROM `debug` WHERE `Clongblob`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `clongblob`=? LIMIT 1
stmt 30 differs:
  orig: SELECT `id` FROM `debug` WHERE `Ctinytext`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `ctinytext`=? LIMIT 1
stmt 31 differs:
  orig: SELECT `id` FROM `debug` WHERE `Ctext`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `ctext`=? LIMIT 1
stmt 32 differs:
  orig: SELECT `id` FROM `debug` WHERE `Cmediumtext`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `cmediumtext`=? LIMIT 1
stmt 33 differs:
  orig: SELECT `id` FROM `debug` WHERE `Clongtext`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `clongtext`=? LIMIT 1
stmt 34 differs:
  orig: SELECT `id` FROM `debug` WHERE `Cenum`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `cenum`=? LIMIT 1
stmt 35 differs:
  orig: SELECT `id` FROM `debug` WHERE `Cset`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `cset`=? LIMIT 1
stmt 36 differs:
  orig: SELECT `id` FROM `debug` WHERE `Cjson`=? LIMIT 1
  fmt:  SELECT `id` FROM `debug` WHERE `cjson`=? LIMIT 1


### coalesce_params/mysql
reason: semantic drift in query.sql
--- query diff:
@@ -1,17 +1,15 @@
 -- name: AddEvent :execlastid
-INSERT INTO `Event` (
-    Timezone
-) VALUES (
-    (CASE WHEN sqlc.arg("Timezone") = "calendar" THEN (SELECT cal.Timezone FROM Calendar cal WHERE cal.IdKey = sqlc.arg("calendarIdKey")) ELSE sqlc.arg("Timezone") END)
-);
+INSERT INTO event (timezone)
+VALUES (CASE WHEN sqlc.arg('Timezone') = 'calendar' THEN (SELECT cal.timezone FROM calendar AS cal WHERE cal.idkey = sqlc.arg('calendarIdKey')) ELSE sqlc.arg('Timezone') END);
 
 -- name: AddAuthor :execlastid
 INSERT INTO authors (
-    address,
-    name,
-    bio
-) VALUES (
-    ?,
-    COALESCE(sqlc.narg("calName"), ""),
-    COALESCE(sqlc.narg("calDescription"), "")
+  address,
+  name,
+  bio
+)
+VALUES (
+  ?,
+  coalesce(sqlc.narg('calName'), ''),
+  coalesce(sqlc.narg('calDescription'), '')
 );
--- semantic check (query.sql):
stmt 1 differs:
  orig: INSERT INTO `Event` (`Timezone`) VALUES (CASE WHEN `sqlc`.`arg`(_UTF8MB4'Timezone')=_UTF8MB4'calendar' THEN (SELECT `cal`.`Timezone` FROM `Calendar` AS `cal` WHERE `cal`.`IdKey`=`sqlc`.`arg`(_UTF8MB4'calendarIdKey')) ELSE `sqlc`.`arg`(_UTF8MB4'Timezone') END)
  fmt:  INSERT INTO `event` (`timezone`) VALUES (CASE WHEN `sqlc`.`arg`(_UTF8MB4'Timezone')=_UTF8MB4'calendar' THEN (SELECT `cal`.`timezone` FROM `calendar` AS `cal` WHERE `cal`.`idkey`=`sqlc`.`arg`(_UTF8MB4'calendarIdKey')) ELSE `sqlc`.`arg`(_UTF8MB4'Timezone') END)


### named_param/pgx/v4
reason: generated code changed beyond SQL text
--- query diff:
@@ -2,20 +2,15 @@
 SELECT name FROM foo WHERE name = sqlc.arg('slug') AND sqlc.arg(filter)::bool;
 
 -- name: AtParams :many
-SELECT name FROM foo WHERE name = @slug AND @filter::bool;
+SELECT name FROM foo WHERE name = @slug AND @ filter::bool;
 
 -- name: InsertFuncParams :one
-INSERT INTO foo(name, bio) values (sqlc.arg('name'), sqlc.arg('bio')) returning name;
+INSERT INTO foo (name, bio) VALUES (sqlc.arg('name'), sqlc.arg('bio')) RETURNING name;
 
 -- name: InsertAtParams :one
-INSERT INTO foo(name, bio) values (@name, @bio) returning name;
-
+INSERT INTO foo (name, bio) VALUES (@name, @bio) RETURNING name;
 
 -- name: Update :one
 UPDATE foo
-SET
-  name = CASE WHEN @set_name::bool
-    THEN @name::text
-    ELSE name
-    END
+SET name = CASE WHEN @ set_name::bool THEN @ name::text ELSE name END
 RETURNING *;
--- structural diff:
--- go/query.sql.go
--- /dev/fd/62	2026-08-28 21:15:41.652347020 +0000
+++ /dev/fd/61	2026-08-28 21:15:41.652347020 +0000
@@ -6,13 +6,8 @@
 
 const atParams = ""
 
-type AtParamsParams struct {
-	Slug	string
-	Filter	bool
-}
-
-func (q *Queries) AtParams(ctx context.Context, arg AtParamsParams) ([]string, error) {
-	rows, err := q.db.Query(ctx, atParams, arg.Slug, arg.Filter)
+func (q *Queries) AtParams(ctx context.Context, slug string) ([]string, error) {
+	rows, err := q.db.Query(ctx, atParams, slug)
 	if err != nil {
 		return nil, err
 	}
@@ -88,13 +83,8 @@
 
 const update = ""
 
-type UpdateParams struct {
-	SetName	bool
-	Name	string
-}
-
-func (q *Queries) Update(ctx context.Context, arg UpdateParams) (Foo, error) {
-	row := q.db.QueryRow(ctx, update, arg.SetName, arg.Name)
+func (q *Queries) Update(ctx context.Context) (Foo, error) {
+	row := q.db.QueryRow(ctx, update)
 	var i Foo
 	err := row.Scan(&i.Name, &i.Bio)
 	return i, err



### named_param/pgx/v5
reason: generated code changed beyond SQL text
--- query diff:
@@ -2,20 +2,15 @@
 SELECT name FROM foo WHERE name = sqlc.arg('slug') AND sqlc.arg(filter)::bool;
 
 -- name: AtParams :many
-SELECT name FROM foo WHERE name = @slug AND @filter::bool;
+SELECT name FROM foo WHERE name = @slug AND @ filter::bool;
 
 -- name: InsertFuncParams :one
-INSERT INTO foo(name, bio) values (sqlc.arg('name'), sqlc.arg('bio')) returning name;
+INSERT INTO foo (name, bio) VALUES (sqlc.arg('name'), sqlc.arg('bio')) RETURNING name;
 
 -- name: InsertAtParams :one
-INSERT INTO foo(name, bio) values (@name, @bio) returning name;
-
+INSERT INTO foo (name, bio) VALUES (@name, @bio) RETURNING name;
 
 -- name: Update :one
 UPDATE foo
-SET
-  name = CASE WHEN @set_name::bool
-    THEN @name::text
-    ELSE name
-    END
+SET name = CASE WHEN @ set_name::bool THEN @ name::text ELSE name END
 RETURNING *;
--- structural diff:
--- go/query.sql.go
--- /dev/fd/62	2026-08-28 21:15:41.692347019 +0000
+++ /dev/fd/61	2026-08-28 21:15:41.692347019 +0000
@@ -6,13 +6,8 @@
 
 const atParams = ""
 
-type AtParamsParams struct {
-	Slug	string
-	Filter	bool
-}
-
-func (q *Queries) AtParams(ctx context.Context, arg AtParamsParams) ([]string, error) {
-	rows, err := q.db.Query(ctx, atParams, arg.Slug, arg.Filter)
+func (q *Queries) AtParams(ctx context.Context, slug string) ([]string, error) {
+	rows, err := q.db.Query(ctx, atParams, slug)
 	if err != nil {
 		return nil, err
 	}
@@ -88,13 +83,8 @@
 
 const update = ""
 
-type UpdateParams struct {
-	SetName	bool
-	Name	string
-}
-
-func (q *Queries) Update(ctx context.Context, arg UpdateParams) (Foo, error) {
-	row := q.db.QueryRow(ctx, update, arg.SetName, arg.Name)
+func (q *Queries) Update(ctx context.Context) (Foo, error) {
+	row := q.db.QueryRow(ctx, update)
 	var i Foo
 	err := row.Scan(&i.Name, &i.Bio)
 	return i, err



### selectstatic/mysql
reason: semantic drift in query.sql; generated code changed beyond SQL text
--- query diff:
@@ -1,2 +1,2 @@
 -- name: SelectStatic :one
-SELECT 'a', 'b' AS b, 1 AS num, true AS truefield, 1.0 AS floater
+SELECT 'a', 'b' AS b, 1 AS num, 1 AS truefield, 0 AS floater;
--- semantic check (query.sql):
stmt 1 differs:
  orig: SELECT _UTF8MB4'a',_UTF8MB4'b' AS `b`,1 AS `num`,TRUE AS `truefield`,1.0 AS `floater`
  fmt:  SELECT _UTF8MB4'a',_UTF8MB4'b' AS `b`,1 AS `num`,1 AS `truefield`,0 AS `floater`
--- structural diff:
--- go/query.sql.go
--- /dev/fd/62	2026-08-28 21:15:41.772347016 +0000
+++ /dev/fd/61	2026-08-28 21:15:41.772347016 +0000
@@ -11,7 +11,7 @@
 	B		string
 	Num		int32
 	Truefield	int32
-	Floater		float64
+	Floater		int32
 }
 
 func (q *Queries) SelectStatic(ctx context.Context) (SelectStaticRow, error) {



### star_expansion_core/mysql
reason: semantic drift in query.sql
--- query diff:
@@ -2,13 +2,13 @@
 SELECT *, *, foo.* FROM foo;
 
 -- name: StarQuotedExpansion :many
-SELECT `t`.* FROM foo `t`;
+SELECT t.* FROM foo AS t;
 
 -- name: StarExpansionJoin :many
-SELECT * FROM foo, bar;
+SELECT * FROM foo JOIN bar;
 
 -- name: StarExpansionSubquery :many
-SELECT * FROM (SELECT * FROM bar) sub;
+SELECT * FROM (SELECT * FROM bar) AS sub;
 
 -- name: StarExpansionCTE :many
 WITH t AS (SELECT * FROM bar) SELECT * FROM t;
--- semantic check (query.sql):
stmt 3 differs:
  orig: SELECT * FROM (`foo`) JOIN `bar`
  fmt:  SELECT * FROM `foo` JOIN `bar`


### unnest/postgresql/pgx/v4
reason: generated code changed beyond SQL text
--- query diff:
@@ -1,9 +1,7 @@
 -- name: CreateMemories :many
 INSERT INTO memories (vampire_id)
-SELECT
-    unnest(@vampire_id::uuid[]) AS vampire_id
-RETURNING
-    *;
+SELECT unnest(@ vampire_id::uuid[]) AS vampire_id
+RETURNING *;
 
 -- name: GetVampireIDs :many
-SELECT vampires.id::uuid FROM unnest(@vampire_id::uuid[]) AS vampires (id);
+SELECT vampires.id::uuid FROM unnest(@ vampire_id::uuid[]) AS vampires(id);
--- structural diff:
--- go/query.sql.go
--- /dev/fd/62	2026-08-28 21:15:41.896347013 +0000
+++ /dev/fd/61	2026-08-28 21:15:41.896347013 +0000
@@ -8,8 +8,8 @@
 
 const createMemories = ""
 
-func (q *Queries) CreateMemories(ctx context.Context, vampireID []uuid.UUID) ([]Memory, error) {
-	rows, err := q.db.Query(ctx, createMemories, vampireID)
+func (q *Queries) CreateMemories(ctx context.Context) ([]Memory, error) {
+	rows, err := q.db.Query(ctx, createMemories)
 	if err != nil {
 		return nil, err
 	}
@@ -35,8 +35,8 @@
 
 const getVampireIDs = ""
 
-func (q *Queries) GetVampireIDs(ctx context.Context, vampireID []uuid.UUID) ([]uuid.UUID, error) {
-	rows, err := q.db.Query(ctx, getVampireIDs, vampireID)
+func (q *Queries) GetVampireIDs(ctx context.Context) ([]uuid.UUID, error) {
+	rows, err := q.db.Query(ctx, getVampireIDs)
 	if err != nil {
 		return nil, err
 	}
--- go/querier.go
--- /dev/fd/62	2026-08-28 21:15:41.912347013 +0000
+++ /dev/fd/61	2026-08-28 21:15:41.912347013 +0000
@@ -7,8 +7,8 @@
 )
 
 type Querier interface {
-	CreateMemories(ctx context.Context, vampireID []uuid.UUID) ([]Memory, error)
-	GetVampireIDs(ctx context.Context, vampireID []uuid.UUID) ([]uuid.UUID, error)
+	CreateMemories(ctx context.Context) ([]Memory, error)
+	GetVampireIDs(ctx context.Context) ([]uuid.UUID, error)
 }
 
 var _ Querier = (*Queries)(nil)



### unnest/postgresql/pgx/v5
reason: generated code changed beyond SQL text
--- query diff:
@@ -1,9 +1,7 @@
 -- name: CreateMemories :many
 INSERT INTO memories (vampire_id)
-SELECT
-    unnest(@vampire_id::uuid[]) AS vampire_id
-RETURNING
-    *;
+SELECT unnest(@ vampire_id::uuid[]) AS vampire_id
+RETURNING *;
 
 -- name: GetVampireIDs :many
-SELECT vampires.id::uuid FROM unnest(@vampire_id::uuid[]) AS vampires (id);
+SELECT vampires.id::uuid FROM unnest(@ vampire_id::uuid[]) AS vampires(id);
--- structural diff:
--- go/query.sql.go
--- /dev/fd/62	2026-08-28 21:15:41.940347012 +0000
+++ /dev/fd/61	2026-08-28 21:15:41.940347012 +0000
@@ -8,8 +8,8 @@
 
 const createMemories = ""
 
-func (q *Queries) CreateMemories(ctx context.Context, vampireID []pgtype.UUID) ([]Memory, error) {
-	rows, err := q.db.Query(ctx, createMemories, vampireID)
+func (q *Queries) CreateMemories(ctx context.Context) ([]Memory, error) {
+	rows, err := q.db.Query(ctx, createMemories)
 	if err != nil {
 		return nil, err
 	}
@@ -35,8 +35,8 @@
 
 const getVampireIDs = ""
 
-func (q *Queries) GetVampireIDs(ctx context.Context, vampireID []pgtype.UUID) ([]pgtype.UUID, error) {
-	rows, err := q.db.Query(ctx, getVampireIDs, vampireID)
+func (q *Queries) GetVampireIDs(ctx context.Context) ([]pgtype.UUID, error) {
+	rows, err := q.db.Query(ctx, getVampireIDs)
 	if err != nil {
 		return nil, err
 	}
--- go/querier.go
--- /dev/fd/62	2026-08-28 21:15:41.956347012 +0000
+++ /dev/fd/61	2026-08-28 21:15:41.956347012 +0000
@@ -7,8 +7,8 @@
 )
 
 type Querier interface {
-	CreateMemories(ctx context.Context, vampireID []pgtype.UUID) ([]Memory, error)
-	GetVampireIDs(ctx context.Context, vampireID []pgtype.UUID) ([]pgtype.UUID, error)
+	CreateMemories(ctx context.Context) ([]Memory, error)
+	GetVampireIDs(ctx context.Context) ([]pgtype.UUID, error)
 }
 
 var _ Querier = (*Queries)(nil)



### unnest/postgresql/stdlib
reason: generated code changed beyond SQL text
--- query diff:
@@ -1,9 +1,7 @@
 -- name: CreateMemories :many
 INSERT INTO memories (vampire_id)
-SELECT
-    unnest(@vampire_id::uuid[]) AS vampire_id
-RETURNING
-    *;
+SELECT unnest(@ vampire_id::uuid[]) AS vampire_id
+RETURNING *;
 
 -- name: GetVampireIDs :many
-SELECT vampires.id::uuid FROM unnest(@vampire_id::uuid[]) AS vampires (id);
+SELECT vampires.id::uuid FROM unnest(@ vampire_id::uuid[]) AS vampires(id);
--- structural diff:
--- go/query.sql.go
--- /dev/fd/62	2026-08-28 21:15:41.984347011 +0000
+++ /dev/fd/61	2026-08-28 21:15:41.984347011 +0000
@@ -4,13 +4,12 @@
 	"context"
 
 	"github.com/google/uuid"
-	"github.com/lib/pq"
 )
 
 const createMemories = ""
 
-func (q *Queries) CreateMemories(ctx context.Context, vampireID []uuid.UUID) ([]Memory, error) {
-	rows, err := q.db.QueryContext(ctx, createMemories, pq.Array(vampireID))
+func (q *Queries) CreateMemories(ctx context.Context) ([]Memory, error) {
+	rows, err := q.db.QueryContext(ctx, createMemories)
 	if err != nil {
 		return nil, err
 	}
@@ -39,8 +38,8 @@
 
 const getVampireIDs = ""
 
-func (q *Queries) GetVampireIDs(ctx context.Context, vampireID []uuid.UUID) ([]uuid.UUID, error) {
-	rows, err := q.db.QueryContext(ctx, getVampireIDs, pq.Array(vampireID))
+func (q *Queries) GetVampireIDs(ctx context.Context) ([]uuid.UUID, error) {
+	rows, err := q.db.QueryContext(ctx, getVampireIDs)
 	if err != nil {
 		return nil, err
 	}
--- go/querier.go
--- /dev/fd/62	2026-08-28 21:15:42.000347011 +0000
+++ /dev/fd/61	2026-08-28 21:15:42.000347011 +0000
@@ -7,8 +7,8 @@
 )
 
 type Querier interface {
-	CreateMemories(ctx context.Context, vampireID []uuid.UUID) ([]Memory, error)
-	GetVampireIDs(ctx context.Context, vampireID []uuid.UUID) ([]uuid.UUID, error)
+	CreateMemories(ctx context.Context) ([]Memory, error)
+	GetVampireIDs(ctx context.Context) ([]uuid.UUID, error)
 }
 
 var _ Querier = (*Queries)(nil)



### unnest_star/postgresql/pgx
reason: generated code changed beyond SQL text
--- query diff:
@@ -1,12 +1,8 @@
 -- name: GetPlanItems :many
 SELECT p.plan_id, p.item_id
-FROM (SELECT * FROM unnest(@ids::bigint[])) AS i(req_item_id),
-LATERAL (
-    SELECT plan_id, item_id
-    FROM plan_items
-    WHERE
-        item_id = i.req_item_id AND
-        (@after = 0 OR plan_id < @after)
-    ORDER BY plan_id DESC
-    LIMIT @limit_count
-) p;
\ No newline at end of file
+FROM (SELECT * FROM unnest(@ ids::bigint[])) AS i(req_item_id), LATERAL (SELECT plan_id, item_id
+FROM plan_items
+WHERE item_id = i.req_item_id
+  AND (@after = 0 OR plan_id < @after)
+ORDER BY plan_id DESC
+LIMIT @limit_count) AS p;
--- structural diff:
--- go/query.sql.go
--- /dev/fd/62	2026-08-28 21:15:42.024347010 +0000
+++ /dev/fd/61	2026-08-28 21:15:42.024347010 +0000
@@ -7,7 +7,6 @@
 const getPlanItems = ""
 
 type GetPlanItemsParams struct {
-	Ids		[]int64
 	After		any
 	LimitCount	int32
 }
@@ -18,7 +17,7 @@
 }
 
 func (q *Queries) GetPlanItems(ctx context.Context, arg GetPlanItemsParams) ([]GetPlanItemsRow, error) {
-	rows, err := q.db.Query(ctx, getPlanItems, arg.Ids, arg.After, arg.LimitCount)
+	rows, err := q.db.Query(ctx, getPlanItems, arg.After, arg.LimitCount)
 	if err != nil {
 		return nil, err
 	}



### datatype/mysql
reason: semantic drift in sql/numeric.sql
--- query diff:
@@ -1,40 +1,7 @@
 -- Numeric Types
 -- https://dev.mysql.com/doc/refman/8.0/en/numeric-type-syntax.html
-CREATE TABLE dt_numeric (
-    a INT,
-    b INTEGER,
-    c TINYINT,
-    d SMALLINT,
-    e MEDIUMINT,
-    f BIGINT,
-    g BIT,
-    h DECIMAL(10, 5),
-    i DEC(10, 5),
-    j FLOAT,
-    k DOUBLE,
-    l DOUBLE PRECISION
-);
+CREATE TABLE dt_numeric (a int, b int, c tinyint, d smallint, e mediumint, f bigint, g bit, h decimal, i decimal, j float, k double, l double);
 
-CREATE TABLE dt_numeric_unsigned (
-    a INT UNSIGNED,
-    b INTEGER UNSIGNED,
-    c TINYINT UNSIGNED,
-    d SMALLINT UNSIGNED,
-    e MEDIUMINT UNSIGNED,
-    f BIGINT UNSIGNED
-);
+CREATE TABLE dt_numeric_unsigned (a int, b int, c tinyint, d smallint, e mediumint, f bigint);
 
-CREATE TABLE dt_numeric_not_null (
-    a INT NOT NULL,
-    b INTEGER NOT NULL,
-    c TINYINT NOT NULL,
-    d SMALLINT NOT NULL,
-    e MEDIUMINT NOT NULL,
-    f BIGINT NOT NULL,
-    g BIT NOT NULL,
-    h DECIMAL(10, 5) NOT NULL,
-    i DEC(10, 5) NOT NULL,
-    j FLOAT NOT NULL,
-    k DOUBLE NOT NULL,
-    l DOUBLE PRECISION NOT NULL
-);
+CREATE TABLE dt_numeric_not_null (a int NOT NULL, b int NOT NULL, c tinyint NOT NULL, d smallint NOT NULL, e mediumint NOT NULL, f bigint NOT NULL, g bit NOT NULL, h decimal NOT NULL, i decimal NOT NULL, j float NOT NULL, k double NOT NULL, l double NOT NULL);
@@ -1,31 +1,5 @@
 -- Character Types
 -- https://dev.mysql.com/doc/refman/8.0/en/string-type-syntax.html
-CREATE TABLE dt_character (
-    a CHARACTER(32),
-    b VARCHAR(32),
-    c CHAR(32),
-    d BINARY(32),
-    e VARBINARY(32),
-    f TINYBLOB,
-    g TINYTEXT,
-    h TEXT,
-    i MEDIUMTEXT,
-    j MEDIUMBLOB,
-    k LONGTEXT,
-    l LONGBLOB
-);
+CREATE TABLE dt_character (a char(32), b varchar(32), c char(32), d binary(32), e varbinary(32), f tinyblob, g tinytext, h text, i mediumtext, j mediumblob, k longtext, l longblob);
 
-CREATE TABLE dt_character_not_null (
-    a CHARACTER(32) NOT NULL,
-    b VARCHAR(32) NOT NULL,
-    c CHAR(32) NOT NULL,
-    d BINARY(32) NOT NULL,
-    e VARBINARY(32) NOT NULL,
-    f TINYBLOB NOT NULL,
-    g TINYTEXT NOT NULL,
-    h TEXT NOT NULL,
-    i MEDIUMTEXT NOT NULL,
-    j MEDIUMBLOB NOT NULL,
-    k LONGTEXT NOT NULL,
-    l LONGBLOB NOT NULL
-);
+CREATE TABLE dt_character_not_null (a char(32) NOT NULL, b varchar(32) NOT NULL, c char(32) NOT NULL, d binary(32) NOT NULL, e varbinary(32) NOT NULL, f tinyblob NOT NULL, g tinytext NOT NULL, h text NOT NULL, i mediumtext NOT NULL, j mediumblob NOT NULL, k longtext NOT NULL, l longblob NOT NULL);
@@ -1,13 +1,5 @@
 -- Date/Time Types
 -- https://www.sqlite.org/datatype3.html
-CREATE TABLE dt_datetime (
-    a DATE,
-    b DATETIME,
-    c TIMESTAMP
-);
+CREATE TABLE dt_datetime (a date, b datetime, c timestamp);
 
-CREATE TABLE dt_datetime_not_null (
-    a DATE NOT NULL,
-    b DATETIME NOT NULL,
-    c TIMESTAMP NOT NULL
-);
+CREATE TABLE dt_datetime_not_null (a date NOT NULL, b datetime NOT NULL, c timestamp NOT NULL);
--- semantic check (sql/numeric.sql):
stmt 1 differs:
  orig: CREATE TABLE `dt_numeric` (`a` INT,`b` INT,`c` TINYINT,`d` SMALLINT,`e` MEDIUMINT,`f` BIGINT,`g` BIT(1),`h` DECIMAL(10,5),`i` DECIMAL(10,5),`j` FLOAT,`k` DOUBLE,`l` DOUBLE)
  fmt:  CREATE TABLE `dt_numeric` (`a` INT,`b` INT,`c` TINYINT,`d` SMALLINT,`e` MEDIUMINT,`f` BIGINT,`g` BIT(1),`h` DECIMAL,`i` DECIMAL,`j` FLOAT,`k` DOUBLE,`l` DOUBLE)
stmt 2 differs:
  orig: CREATE TABLE `dt_numeric_unsigned` (`a` INT UNSIGNED,`b` INT UNSIGNED,`c` TINYINT UNSIGNED,`d` SMALLINT UNSIGNED,`e` MEDIUMINT UNSIGNED,`f` BIGINT UNSIGNED)
  fmt:  CREATE TABLE `dt_numeric_unsigned` (`a` INT,`b` INT,`c` TINYINT,`d` SMALLINT,`e` MEDIUMINT,`f` BIGINT)
stmt 3 differs:
  orig: CREATE TABLE `dt_numeric_not_null` (`a` INT NOT NULL,`b` INT NOT NULL,`c` TINYINT NOT NULL,`d` SMALLINT NOT NULL,`e` MEDIUMINT NOT NULL,`f` BIGINT NOT NULL,`g` BIT(1) NOT NULL,`h` DECIMAL(10,5) NOT NULL,`i` DECIMAL(10,5) NOT NULL,`j` FLOAT NOT NULL,`k` DOUBLE NOT NULL,`l` DOUBLE NOT NULL)
  fmt:  CREATE TABLE `dt_numeric_not_null` (`a` INT NOT NULL,`b` INT NOT NULL,`c` TINYINT NOT NULL,`d` SMALLINT NOT NULL,`e` MEDIUMINT NOT NULL,`f` BIGINT NOT NULL,`g` BIT(1) NOT NULL,`h` DECIMAL NOT NULL,`i` DECIMAL NOT NULL,`j` FLOAT NOT NULL,`k` DOUBLE NOT NULL,`l` DOUBLE NOT NULL)


### mysql_reference_manual
reason: semantic drift in aggregate_functions/group_concat.sql
--- query diff:
@@ -1,9 +1,8 @@
 -- name: DateSubOneYear :one
-SELECT DATE_SUB('2018-05-01',INTERVAL 1 YEAR);
+SELECT date_sub('2018-05-01', INTERVAL 1 YEAR);
 
 -- name: DateSubDaySecond :one
-SELECT DATE_SUB('2025-01-01 00:00:00',
-                INTERVAL '1 1:1:1' DAY_SECOND);
+SELECT date_sub('2025-01-01 00:00:00', INTERVAL '1 1:1:1' DAY_SECOND);
 
 -- name: DateSub31Days :one
-SELECT DATE_SUB('1998-01-02', INTERVAL 31 DAY);
+SELECT date_sub('1998-01-02', INTERVAL 31 DAY);
@@ -1,24 +1,19 @@
 -- https://dev.mysql.com/doc/refman/8.0/en/date-and-time-functions.html#function_date-add
 
 -- name: DateAddOneDay :one
-SELECT DATE_ADD('2018-05-01',INTERVAL 1 DAY);
+SELECT date_add('2018-05-01', INTERVAL 1 DAY);
 
 -- name: DateAddOneSecond :one
-SELECT DATE_ADD('2020-12-31 23:59:59',
-                INTERVAL 1 SECOND);
+SELECT date_add('2020-12-31 23:59:59', INTERVAL 1 SECOND);
 
 -- name: DateAddTimestampOneSecond :one
-SELECT DATE_ADD('2018-12-31 23:59:59',
-                INTERVAL 1 DAY);
+SELECT date_add('2018-12-31 23:59:59', INTERVAL 1 DAY);
 
 -- name: DateAddMinuteSecond :one
-SELECT DATE_ADD('2100-12-31 23:59:59',
-                INTERVAL '1:1' MINUTE_SECOND);
+SELECT date_add('2100-12-31 23:59:59', INTERVAL '1:1' MINUTE_SECOND);
 
 -- name: DateAddDayHour :one
-SELECT DATE_ADD('1900-01-01 00:00:00',
-                INTERVAL '-1 10' DAY_HOUR);
+SELECT date_add('1900-01-01 00:00:00', INTERVAL '-1 10' DAY_HOUR);
 
 -- name: DateAddSecondMicrosecond :one
-SELECT DATE_ADD('1992-12-31 23:59:59.000002',
-           INTERVAL '1.999999' SECOND_MICROSECOND);
+SELECT date_add('1992-12-31 23:59:59.000002', INTERVAL '1.999999' SECOND_MICROSECOND);
@@ -1,10 +1,11 @@
 -- name: GroupConcat :many
-SELECT student_name, GROUP_CONCAT(test_score)
+SELECT student_name, group_concat(test_score)
 FROM student
 GROUP BY student_name;
 
 -- name: GroupConcatOrderBy :many
-SELECT student_name,
-    GROUP_CONCAT(DISTINCT test_score ORDER BY test_score DESC SEPARATOR ' ')
+SELECT
+  student_name,
+  group_concat(DISTINCT test_score SEPARATOR ' ')
 FROM student
 GROUP BY student_name;
--- semantic check (aggregate_functions/group_concat.sql):
stmt 2 differs:
  orig: SELECT `student_name`,GROUP_CONCAT(DISTINCT `test_score` ORDER BY `test_score` DESC SEPARATOR ' ') FROM `student` GROUP BY `student_name`
  fmt:  SELECT `student_name`,GROUP_CONCAT(DISTINCT `test_score` SEPARATOR ' ') FROM `student` GROUP BY `student_name`


