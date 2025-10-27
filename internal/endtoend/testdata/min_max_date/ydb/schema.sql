CREATE TABLE activities (
    account_id Int64 NOT NULL,
    event_time Timestamp NOT NULL,
    PRIMARY KEY (account_id, event_time)
);

