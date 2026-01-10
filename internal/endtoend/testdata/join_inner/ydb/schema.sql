CREATE TABLE events (
    ID Int32 NOT NULL,
    PRIMARY KEY (ID)
);

CREATE TABLE handled_events (
    last_handled_id Int32,
    handler Utf8,
    PRIMARY KEY (last_handled_id, handler)
);




