CREATE SCHEMA foo;

CREATE TYPE point_type AS (
    x integer,
    y integer
);

CREATE TYPE foo.point_type AS (
    x integer,
    y integer
);

CREATE TABLE foo.paths (
    point_one point_type,
    point_two foo.point_type
);

