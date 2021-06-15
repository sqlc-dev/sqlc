CREATE TABLE foo (val text not null);

CREATE VIEW first_view AS SELECT * FROM foo;
CREATE VIEW second_view AS SELECT * FROM foo;
CREATE VIEW third_view AS SELECT * FROM foo;

ALTER TABLE foo ADD COLUMN val2 integer;
CREATE OR REPLACE VIEW second_view AS SELECT * FROM foo;

DROP VIEW third_view;
