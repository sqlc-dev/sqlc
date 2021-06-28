CREATE TABLE foo (
    foobar ENUM ('foo-a', 'foo_b', 'foo:c', 'foo/d', 'foo@e', 'foo+f', 'foo!g') NOT NULL,
    digit  ENUM ('0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '#', '*')    NOT NULL
);
