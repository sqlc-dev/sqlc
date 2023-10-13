CREATE EXTENSION IF NOT EXISTS ltree;

CREATE TABLE foo (
    qualified_name ltree,
    name_query lquery,
    fts_name_query ltxtquery
);

