CREATE TABLE metrics (
    id      INT          NOT NULL PRIMARY KEY,
    value   FLOAT        NULL,
    count   INT          NOT NULL,
    ratio   DOUBLE       NULL
);
