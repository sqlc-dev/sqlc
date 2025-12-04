CREATE TABLE tbl (
        value integer
);

-- https://www.postgresql.org/docs/current/sql-createprocedure.html
CREATE PROCEDURE insert_data(IN a integer, IN b integer, OUT c integer)
    LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO tbl VALUES (a);
    INSERT INTO tbl VALUES (b);

    c := 777;
END;
$$;
