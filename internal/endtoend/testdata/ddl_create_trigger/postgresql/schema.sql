CREATE SCHEMA tdd;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto
    SCHEMA public
    VERSION "1.3";

CREATE FUNCTION tdd.trigger_set_timestamp() RETURNS trigger
LANGUAGE plpgsql
AS $$BEGIN
  NEW.ts_updated = NOW();
  RETURN NEW;
END;
$$;

CREATE TABLE tdd.tests (
    test_id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    title text DEFAULT ''::text NOT NULL,
    descr text DEFAULT ''::text NOT NULL,
    ts_created timestamp with time zone DEFAULT now() NOT NULL,
    ts_updated timestamp with time zone DEFAULT now() NOT NULL
);

CREATE TRIGGER set_timestamp BEFORE UPDATE ON tdd.tests FOR EACH ROW EXECUTE FUNCTION tdd.trigger_set_timestamp();
