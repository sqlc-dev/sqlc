ALTER TABLE applications ADD COLUMN build_nupkg BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE applications ADD COLUMN allow_build_nupkg BOOLEAN NOT NULL DEFAULT false;
