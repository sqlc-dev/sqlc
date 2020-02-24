CREATE TABLE foo (
        bar text NOT NULL 
);

CREATE TABLE bar (
        baz text NOT NULL 
);

SELECT bar FROM foo;

DROP TABLE bar;
DROP TABLE IF EXISTS baz;
