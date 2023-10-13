-- FILE: schema.sql

DROP SCHEMA IF EXISTS td3 CASCADE;
CREATE SCHEMA td3;

CREATE TABLE td3.codes (
            id SERIAL PRIMARY KEY,
            ts_created timestamptz DEFAULT now() NOT NULL,
            ts_updated timestamptz DEFAULT now() NOT NULL,
            created_by text NOT NULL,
            updated_by text NOT NULL,
            
            code text,
            hash text,
            is_private boolean
);


CREATE TABLE td3.test_codes (
            id SERIAL PRIMARY KEY,
            ts_created timestamptz DEFAULT now() NOT NULL,
            ts_updated timestamptz DEFAULT now() NOT NULL,
            created_by text NOT NULL,
            updated_by text NOT NULL,

            test_id integer NOT NULL,
            code_hash text NOT NULL
);

-- FILE: query.sql

