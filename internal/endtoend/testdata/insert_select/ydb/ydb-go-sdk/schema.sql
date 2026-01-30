CREATE TABLE bar (
    name Utf8 NOT NULL,
    ready Bool NOT NULL,
    PRIMARY KEY (name)
);

CREATE TABLE foo (
    name Utf8 NOT NULL,
    meta Utf8 NOT NULL,
    PRIMARY KEY (name)
);
