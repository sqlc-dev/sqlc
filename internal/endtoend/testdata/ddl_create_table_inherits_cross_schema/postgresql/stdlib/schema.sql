CREATE SCHEMA parent;
CREATE SCHEMA child;

CREATE TABLE parent.party (
    party_id uuid PRIMARY KEY,
    name text NOT NULL
);

CREATE TABLE child.person (
    first_name text NOT NULL,
    last_name text NOT NULL
) INHERITS (parent.party);

