CREATE TABLE user_referrals
(
    user_id     bigint                                 NOT NULL
        CONSTRAINT user_referrals_pk
            PRIMARY KEY,
    referrer_id bigint                                 NOT NULL,
    balance     numeric(38, 9)           DEFAULT 0.00  NOT NULL,
    created_at  timestamp with time zone DEFAULT NOW() NOT NULL,
    updated_at  timestamp with time zone DEFAULT NOW() NOT NULL
);