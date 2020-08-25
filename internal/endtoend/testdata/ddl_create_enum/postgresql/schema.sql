CREATE TYPE foobar AS ENUM (
    -- Valid separators
    'foo-a',
    'foo_b',
    'foo:c',
    'foo/d',
    -- Strip unknown characters
    'foo@e',
    'foo+f',
    'foo!g'
);

CREATE TABLE foo (val foobar NOT NULL);
