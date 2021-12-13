CREATE TABLE party (
    party_id uuid PRIMARY KEY,
    name text NOT NULL
);

CREATE TABLE person (
    first_name text NOT NULL,
    last_name text NOT NULL
) INHERITS (party);

CREATE TABLE organisation (
    legal_name text
) INHERITS (party);
