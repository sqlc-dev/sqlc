CREATE TABLE bar (
    id integer not null,
    UNIQUE(id)
);

CREATE TABLE foo (id integer not null, bar_id integer references bar(id));


