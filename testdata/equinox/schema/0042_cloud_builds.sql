CREATE TABLE cloud_builds (
        appsid    varchar(64)  NOT NULL REFERENCES applications(sid),
        finished  timestamp    DEFAULT NULL,
        created   timestamp    DEFAULT NOW(),
        status    text,
        id        text
);

-- Ensure that we can only create one cloud build at a time
CREATE UNIQUE INDEX cloud_builds_constraint ON cloud_builds(appsid)
    WHERE finished IS NULL;
