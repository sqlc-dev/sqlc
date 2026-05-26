CREATE TABLE not_null_leg (
    id          bigint  NOT NULL,
    label       text    NOT NULL,
    rank        integer NOT NULL
);

CREATE TABLE nullable_leg (
    id          bigint,
    label       text,
    rank        integer NOT NULL
);
