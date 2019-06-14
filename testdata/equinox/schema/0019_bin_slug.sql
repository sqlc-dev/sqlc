ALTER TABLE applications ADD COLUMN bin_slug VARCHAR(255);
UPDATE applications set bin_slug = slug;
ALTER TABLE applications ALTER COLUMN bin_slug SET NOT NULL;
