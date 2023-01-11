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

-- name: InsertCode :one
WITH cc AS (
            INSERT INTO td3.codes(created_by, updated_by, code, hash, is_private)
            VALUES ($1, $1, $2, $3, false)
            RETURNING hash
)
INSERT INTO td3.test_codes(created_by, updated_by, test_id, code_hash)
VALUES(
            $1, $1, $4, (select hash from cc)
)
RETURNING *;
