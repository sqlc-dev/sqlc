CREATE TABLE tbl (
        value integer
);

-- https://www.postgresql.org/docs/current/sql-createprocedure.html
CREATE PROCEDURE insert_data(a integer, b integer)
LANGUAGE SQL
AS $$
INSERT INTO tbl VALUES (a);
INSERT INTO tbl VALUES (b);
$$;
