CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE foo (
    bar citext,
    bat citext not null
);

