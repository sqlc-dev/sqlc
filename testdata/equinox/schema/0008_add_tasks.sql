CREATE TYPE runstate AS ENUM ('started', 'finished');

CREATE TABLE tasks (
        id             SERIAL       UNIQUE NOT NULL,
        sid            varchar(64)  UNIQUE NOT NULL,
        state          runstate     DEFAULT 'started',
        buildid        bigint       NOT NULL REFERENCES builds(id),
        rawassetid     bigint       NOT NULL REFERENCES assets(id),
        target         varchar(64)  NOT NULL,
        assetsid       varchar(64)  UNIQUE NOT NULL,
        assetchecksum  TEXT         DEFAULT '',
        assetsize      int          DEFAULT 0,
        created        timestamp    DEFAULT NOW(),
        finished       timestamp    DEFAULT NULL
);

ALTER TABLE builds ADD COLUMN state runstate DEFAULT 'started';
ALTER TABLE builds DROP COLUMN started;
UPDATE builds SET state = 'finished' WHERE finished IS NOT NULL;
