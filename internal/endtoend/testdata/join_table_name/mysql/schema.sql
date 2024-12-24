CREATE TABLE bar (
    id integer not null,
    UNIQUE (id)
);

CREATE TABLE foo (id integer not null, bar integer references bar(id));

