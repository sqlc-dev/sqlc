CREATE TABLE foo (
    other     Text NOT NULL,
    total     Int64 NOT NULL,
    tags      Text NOT NULL,
    byte_seq  String NOT NULL,
    retyped   Text NOT NULL,
    langs     Text,
    PRIMARY KEY (other, total)
);
