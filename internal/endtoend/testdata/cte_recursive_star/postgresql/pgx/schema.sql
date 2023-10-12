create table dict(
    id           VARCHAR(36)     PRIMARY KEY DEFAULT gen_random_uuid(),
    app_id       VARCHAR(36)     NOT NULL,
    code         VARCHAR(64),
    parent_code  VARCHAR(64)     NOT NULL,
    label        TEXT            NOT NULL DEFAULT '',
    value        TEXT            NULL,
    weight       INT             NOT NULL DEFAULT 0,
    is_default   BOOLEAN         NOT NULL DEFAULT false,
    is_virtual   BOOLEAN         NOT NULL DEFAULT false,
    status       SMALLINT        NOT NULL DEFAULT 1,
    create_at    TIMESTAMPTZ(0)  NOT NULL DEFAULT now(),
    create_by    VARCHAR(36)     NOT NULL DEFAULT '',
    update_at    TIMESTAMPTZ(0),
    update_by    VARCHAR(36),
    is_delete    BOOLEAN         NOT NULL DEFAULT false
);
