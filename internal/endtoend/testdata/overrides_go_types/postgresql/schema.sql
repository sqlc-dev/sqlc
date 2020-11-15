CREATE TABLE new_style (
    id         uuid NOT NULL,
    other_id   uuid NOT NULL,
    age        integer,
    balance    double,
    bio        text,
    about      text
);

CREATE TABLE old_style (
    id         uuid NOT NULL,
    about      text
);
