CREATE TABLE tbl (
    value int
);

CREATE PROCEDURE insert_data(a int, b int)
BEGIN
    INSERT INTO tbl VALUES (a);
    INSERT INTO tbl VALUES (b);
END;