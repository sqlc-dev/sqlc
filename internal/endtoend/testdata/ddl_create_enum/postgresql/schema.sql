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

CREATE TYPE "digit" AS ENUM (
  '0',
  '1',
  '2',
  '3',
  '4',
  '5',
  '6',
  '7',
  '8',
  '9',
  '#',
  '*'
);

CREATE TABLE foo (val foobar NOT NULL);
