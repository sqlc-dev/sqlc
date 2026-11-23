CREATE TABLE plan (
    id UUID NOT NULL,
    ship_size VARCHAR NOT NULL,
    price DOUBLE
);

-- A declared column takes the place of the pseudo-column.
CREATE TABLE ledger (
    rowid VARCHAR NOT NULL,
    amount INTEGER
);
