ALTER TABLE account ADD COLUMN stripe_sub text;
UPDATE account
SET stripe_sub = subscription.stripe_id
FROM subscription 
WHERE account.id = subscription.account_id;
