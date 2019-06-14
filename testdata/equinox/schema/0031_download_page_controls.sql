ALTER table applications ADD COLUMN download_page_username VARCHAR(255) NOT NULL default '';
ALTER table applications ADD COLUMN download_page_password VARCHAR(255) NOT NULL default '';
ALTER table applications ADD COLUMN allow_download_page_controls BOOL DEFAULT false;
UPDATE applications set allow_download_page_controls = (SELECT plan NOT IN ('basic_yearly', 'basic_monthly') FROM account WHERE account.id = applications.accountid);
