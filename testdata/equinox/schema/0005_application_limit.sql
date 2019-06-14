-- The accounts table is owned be the accounts service, so don't add
-- columsn to it.
CREATE TABLE application_limits (
  accountid   bigint NOT NULL,
  total_apps  bigint DEFAULT 0,
  max_apps    bigint DEFAULT 10,

  CONSTRAINT account_application_limit CHECK (total_apps <= max_apps)
);

INSERT INTO application_limits (accountid) SELECT id FROM account;
UPDATE application_limits
SET total_apps = (
  SELECT COUNT(1) FROM applications
  WHERE applications.accountid = application_limits.accountid
);

CREATE OR REPLACE FUNCTION increment_app_count() RETURNS trigger AS $$
DECLARE
BEGIN
    UPDATE application_limits SET total_apps = total_apps+1 WHERE accountid = NEW.accountid;
    IF NOT FOUND THEN
    INSERT INTO application_limits (accountid, total_apps) VALUES (NEW.accountid, 1);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

CREATE TRIGGER increment_app
    BEFORE INSERT ON applications
    FOR EACH ROW
    EXECUTE PROCEDURE increment_app_count();

CREATE OR REPLACE FUNCTION decrement_app_count() RETURNS trigger AS $$
DECLARE
BEGIN
    UPDATE application_limits SET total_apps = total_apps-1 WHERE accountid = OLD.accountid;
    RETURN OLD;
END;
$$ LANGUAGE 'plpgsql';

CREATE TRIGGER decrement_app
    BEFORE DELETE ON applications
    FOR EACH ROW
    EXECUTE PROCEDURE decrement_app_count();

