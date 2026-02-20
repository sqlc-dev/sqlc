CREATE TYPE party_rank AS ENUM('ensign', 'lieutenant', 'commander');

CREATE TABLE party (
    party_id uuid PRIMARY KEY,
    name text NOT NULL,
    joined_on timestamptz NOT NULL,
    ronk party_rank,
    unnecessary_column text
);

CREATE TABLE person (
    first_name text NOT NULL,
    last_name text NOT NULL
) INHERITS (party);

CREATE TABLE organisation (
    legal_name text
) INHERITS (party);

CREATE TABLE llc (
   incorporation_date timestamp,
   legal_name text NOT NULL
) INHERITS (organisation);


ALTER TABLE party ALTER COLUMN ronk SET NOT NULL; 
ALTER TABLE party RENAME COLUMN ronk TO rank;
ALTER TABLE party ALTER COLUMN joined_on DROP NOT NULL;
ALTER TABLE party DROP COLUMN unnecessary_column;
ALTER TABLE party ADD COLUMN region text;

ALTER TYPE party_rank ADD VALUE 'captain';
ALTER TYPE party_rank ADD VALUE 'admiral';
