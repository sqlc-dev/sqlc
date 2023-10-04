CREATE TABLE authors1 (
  id   INTEGER PRIMARY KEY,
  name text      NOT NULL,
  bio  text
) STRICT, WITHOUT ROWID;

CREATE TABLE authors2 (
  id   INTEGER PRIMARY KEY,
  name text      NOT NULL,
  bio  text
) WITHOUT ROWID, STRICT;

CREATE TABLE authors3 (
  id   INTEGER PRIMARY KEY,
  name text      NOT NULL,
  bio  text
) WITHOUT ROWID;

CREATE TABLE authors4 (
  id   INTEGER PRIMARY KEY,
  name text      NOT NULL,
  bio  text
) STRICT;

