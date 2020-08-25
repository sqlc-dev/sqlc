CREATE TABLE foo (
    foobar ENUM ('foo-a', 'foo_b', 'foo:c', 'foo/d', 'foo@e', 'foo+f', 'foo!g') NOT NULL
);
