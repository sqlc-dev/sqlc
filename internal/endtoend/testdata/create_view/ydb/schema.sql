CREATE TABLE foo (
    val Text NOT NULL,
    PRIMARY KEY (val)
);

CREATE VIEW first_view AS SELECT * FROM foo;
CREATE VIEW third_view AS SELECT * FROM foo;

ALTER TABLE foo ADD COLUMN val2 Int32;
-- YDB doesn't support CREATE OR REPLACE VIEW, only CREATE or DROP
-- So we need to DROP and CREATE again
DROP VIEW second_view;
CREATE VIEW second_view AS SELECT * FROM foo;

DROP VIEW third_view;
