CREATE TABLE applications (
        id         SERIAL       UNIQUE NOT NULL,
        sid        varchar(64)  UNIQUE NOT NULL,
        accountid  bigint       NOT NULL,
        created    timestamp    DEFAULT NOW(),
        name       varchar(255) NOT NULL,

        CONSTRAINT account_id_name UNIQUE(accountid, name)
);

CREATE TABLE credentials (
        id         SERIAL       UNIQUE NOT NULL,
        sid        varchar(64)  UNIQUE NOT NULL,
        created    timestamp    DEFAULT NOW(),
        accountid  bigint       NOT NULL,
        tokenhash  varchar(255) NOT NULL
);

CREATE TYPE release_state AS ENUM ('draft','building','published','revoked');

CREATE TABLE releases (
        id             SERIAL         UNIQUE NOT NULL,
        sid            varchar(64)    UNIQUE NOT NULL,
        created        timestamp      DEFAULT NOW(),
        appid          bigint         NOT NULL REFERENCES applications(id),
        version        varchar(64)    NOT NULL,
        state          release_state  NOT NULL,
        title          varchar(255)   DEFAULT '',
        description    text           DEFAULT '',

        CONSTRAINT     app_id_version UNIQUE(appid, version)
);

CREATE TABLE assets (
        id             SERIAL       UNIQUE NOT NULL,
        sid            varchar(64)  UNIQUE NOT NULL,
        created        timestamp    DEFAULT NOW(),
        releaseid      bigint       NOT NULL REFERENCES releases(id),
        os             varchar(64)  NOT NULL,
        arch           varchar(64)  NOT NULL,
        archiveformat  varchar(64)  NOT NULL,
        goarm          varchar(64)  DEFAULT '',
        checksum       varchar(64)  DEFAULT '',
        signature      varchar(255) DEFAULT '',
        compiler       jsonb DEFAULT '{}',
        size           bigint NOT NULL
);

CREATE TABLE channels (
        id             SERIAL       UNIQUE NOT NULL,
        sid            varchar(64)  UNIQUE NOT NULL,
        created        timestamp    DEFAULT NOW(),
        appid          bigint       NOT NULL REFERENCES applications(id),
        name           varchar(255) NOT NULL,
        title          varchar(255) DEFAULT '',
        description    text         DEFAULT '',

        CONSTRAINT     app_id_name  UNIQUE(appid, name)
);

CREATE TABLE publishings (
        id             SERIAL       UNIQUE NOT NULL,
        sid            varchar(64)  UNIQUE NOT NULL,
        channelid      bigint       NOT NULL REFERENCES channels(id),
        releaseid      bigint       NOT NULL REFERENCES releases(id),
        created        timestamp    DEFAULT NOW(),

        CONSTRAINT     channel_id_release_id  UNIQUE(channelid, releaseid)
);

CREATE TABLE builds (
        id             SERIAL       UNIQUE NOT NULL,
        sid            varchar(64)  UNIQUE NOT NULL,
        releaseid      bigint       NOT NULL REFERENCES releases(id),
        channelid      bigint       REFERENCES channels(id),
        started        timestamp    DEFAULT NULL,
        finished       timestamp    DEFAULT NULL,
        created        timestamp    DEFAULT NOW()
);
