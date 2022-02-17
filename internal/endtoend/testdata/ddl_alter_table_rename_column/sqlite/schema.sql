CREATE TABLE foo (bar text);
ALTER TABLE foo RENAME COLUMN bar TO baz;

ALTER TABLE foo RENAME baz TO boo;
