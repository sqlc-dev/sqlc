-- https://github.com/sqlc-dev/sqlc/issues/1728

CREATE TABLE locations (
    id              Serial,
    name            Text NOT NULL,
    address         Text NOT NULL,
    zip_code        Int32 NOT NULL,
    latitude        Double NOT NULL,
    longitude       Double NOT NULL,
    PRIMARY KEY (id)
);



