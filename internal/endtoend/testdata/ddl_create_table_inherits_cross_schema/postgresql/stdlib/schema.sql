CREATE SCHEMA one;
CREATE SCHEMA two;

CREATE TABLE one.party (
    party_id uuid PRIMARY KEY,
    name text NOT NULL
);

CREATE TABLE two.person (
    first_name text NOT NULL,
    last_name text NOT NULL
) INHERITS (one.party);

