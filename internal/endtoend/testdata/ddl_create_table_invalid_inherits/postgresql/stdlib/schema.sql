CREATE TABLE party (
    name text NOT NULL
);

CREATE TABLE organisation (
    name integer NOT NULL
) INHERITS (party);
