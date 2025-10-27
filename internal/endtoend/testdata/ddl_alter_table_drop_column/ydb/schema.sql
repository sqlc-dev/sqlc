CREATE TABLE foo (
    bar Text,
    baz Text,
    PRIMARY KEY (bar)
);

ALTER TABLE foo DROP COLUMN baz;
