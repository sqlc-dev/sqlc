ALTER TABLE applications ADD COLUMN slug VARCHAR(255);
ALTER TABLE applications ADD CONSTRAINT account_id_slug UNIQUE (accountid, slug);
UPDATE applications SET slug=name;
ALTER TABLE applications ALTER COLUMN slug SET NOT NULL;
ALTER TABLE applications DROP CONSTRAINT account_id_name;
