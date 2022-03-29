CREATE TABLE foo (val text not null);

CREATE MATERIALIZED VIEW mat_first_view AS SELECT * FROM foo;
