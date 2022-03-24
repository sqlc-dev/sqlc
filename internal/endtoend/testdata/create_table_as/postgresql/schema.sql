CREATE TABLE foo (val text not null);

CREATE TABLE second_table AS SELECT * FROM foo;
