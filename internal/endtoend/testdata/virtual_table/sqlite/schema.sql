CREATE TABLE tbl(a INTEGER PRIMARY KEY, b TEXT, c TEXT, d TEXT, e INTEGER);

CREATE VIRTUAL TABLE tbl_ft USING fts5(b, c UNINDEXED, content='tbl', content_rowid='a');

CREATE VIRTUAL TABLE ft USING fts5(b);

CREATE TRIGGER tbl_ai AFTER INSERT ON tbl BEGIN
  INSERT INTO tbl_ft(rowid, b, c) VALUES (new.a, new.b, new.c);
END;

INSERT INTO tbl VALUES(1, 'xx yy cc', 't', 'a', 11);
INSERT INTO tbl VALUES(2, 'aa bb', 't', 'a', 22);

INSERT INTO ft VALUES('xx cc');
INSERT INTO ft VALUES('cc bb');
