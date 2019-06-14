ALTER table account_limits ADD COLUMN allow_download_page_controls BOOL DEFAULT false;
UPDATE account_limits set allow_download_page_controls = (SELECT plan NOT IN ('basic_yearly', 'basic_monthly') FROM account WHERE account.id = account_limits.accountid);

