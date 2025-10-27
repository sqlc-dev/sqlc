CREATE TABLE bar (
    id Int32 NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE foo (
    id Int32 NOT NULL,
    bar_id Int32,
    PRIMARY KEY (id)
);


