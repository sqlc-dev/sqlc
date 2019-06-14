DROP TRIGGER increment_app on applications;
DROP TRIGGER decrement_app on applications;
DROP FUNCTION decrement_app_count();
DROP FUNCTION increment_app_count();

ALTER TABLE application_limits RENAME TO account_limits;
ALTER TABLE account_limits ADD COLUMN max_releases bigint DEFAULT 1000;
ALTER TABLE account_limits ADD COLUMN max_channels bigint DEFAULT 10;
ALTER TABLE account_limits ADD CONSTRAINT accountid_key UNIQUE (accountid);

CREATE OR REPLACE FUNCTION increment_app_count() RETURNS trigger AS $$
DECLARE
BEGIN
    UPDATE account_limits SET total_apps = total_apps+1 WHERE accountid = NEW.accountid;
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
    UPDATE account_limits SET total_apps = total_apps-1 WHERE accountid = OLD.accountid;
    RETURN OLD;
END;
$$ LANGUAGE 'plpgsql';

CREATE TRIGGER decrement_app
    BEFORE DELETE ON applications
    FOR EACH ROW
    EXECUTE PROCEDURE decrement_app_count();
