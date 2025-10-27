CREATE TABLE foo (
    bar Text,
    PRIMARY KEY (bar)
);

-- YDB doesn't support ALTER COLUMN TYPE, so we use DROP COLUMN + ADD COLUMN
ALTER TABLE foo DROP COLUMN bar;
ALTER TABLE foo ADD COLUMN bar Timestamp;
