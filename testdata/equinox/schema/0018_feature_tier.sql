ALTER TABLE applications ADD COLUMN allow_build_msi bool NOT NULL DEFAULT true;
ALTER TABLE applications ADD COLUMN allow_build_pkg bool NOT NULL DEFAULT true;
ALTER TABLE applications ADD COLUMN allow_build_rpm bool NOT NULL DEFAULT true;
ALTER TABLE applications ADD COLUMN allow_build_deb bool NOT NULL DEFAULT true;
ALTER TABLE applications ADD COLUMN allow_homebrew bool NOT NULL DEFAULT true;

ALTER TABLE account_limits ADD COLUMN allow_build_msi bool NOT NULL DEFAULT true;
ALTER TABLE account_limits ADD COLUMN allow_build_pkg bool NOT NULL DEFAULT true;
ALTER TABLE account_limits ADD COLUMN allow_build_rpm bool NOT NULL DEFAULT true;
ALTER TABLE account_limits ADD COLUMN allow_build_deb bool NOT NULL DEFAULT true;
ALTER TABLE account_limits ADD COLUMN allow_homebrew bool NOT NULL DEFAULT true;
