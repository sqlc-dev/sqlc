ALTER TABLE applications ADD COLUMN product varchar(64);
UPDATE applications SET product = CASE 
  WHEN account.plan = 'basic_yearly' THEN 'basic'
  WHEN account.plan = 'basic_monthly' THEN 'basic'
  WHEN account.plan = 'business_monthly' THEN 'business'
  WHEN account.plan = 'business_yearly' THEN 'business'
  ELSE 'free'
END FROM account WHERE account.id = applications.accountid;
ALTER TABLE applications ALTER COLUMN product SET NOT NULL;
