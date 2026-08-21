CREATE TABLE foo (bar text, baz text);
ALTER TABLE foo
    RENAME COLUMN bar TO qux,
    RENAME COLUMN baz TO quux;
