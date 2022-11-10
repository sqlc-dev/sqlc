CREATE TABLE foo (bar text NOT NULL, baz text NOT NULL);
ALTER TABLE foo MODIFY COLUMN bar integer;
ALTER TABLE foo MODIFY baz integer;
