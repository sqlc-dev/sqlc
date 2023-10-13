CREATE TABLE authors (
  id          BIGINT PRIMARY KEY,
  foo         ENUM("ok") DEFAULT "ok" NOT NULL,
  renamed     ENUM("ok") DEFAULT "ok" NOT NULL,
  removed     ENUM("ok") DEFAULT "ok" NOT NULL,
  add_item    ENUM("ok") DEFAULT "ok" NOT NULL,
  remove_item ENUM("ok", "removed") DEFAULT "ok" NOT NULL
);

CREATE TABLE renamed (
  id       BIGINT PRIMARY KEY,
  foo      ENUM("ok") DEFAULT "ok" NOT NULL
);

CREATE TABLE removed (
  id       BIGINT PRIMARY KEY,
  foo      ENUM("ok") DEFAULT "ok" NOT NULL
);

/* Rename column */
ALTER TABLE authors RENAME COLUMN renamed TO bar;

/* Drop column */
ALTER TABLE authors DROP COLUMN removed;

/* Add column */
ALTER TABLE authors ADD COLUMN added ENUM("ok") DEFAULT "ok" NOT NULL;

/* Add enum values */
ALTER TABLE authors MODIFY add_item ENUM("ok", "added") DEFAULT "ok" NOT NULL;

/* Remove enum values */
ALTER TABLE authors MODIFY remove_item ENUM("ok") DEFAULT "ok" NOT NULL;

/* Drop table */
DROP TABLE removed;

/* Rename table */
ALTER TABLE renamed RENAME TO books;
