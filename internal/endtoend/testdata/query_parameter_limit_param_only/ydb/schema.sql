CREATE TABLE notice (
    id Int32 NOT NULL,
    cnt Int32 NOT NULL,
    status Text NOT NULL,
    notice_at Timestamp,
    created_at Timestamp NOT NULL,
    PRIMARY KEY (id)
);
