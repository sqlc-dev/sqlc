CREATE TABLE foo (
	PRIMARY KEY (a, b) INCLUDE (c),
	a integer,
	b integer,
	c integer
);
