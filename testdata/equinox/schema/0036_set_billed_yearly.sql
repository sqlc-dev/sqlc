UPDATE account SET billed_yearly = true
WHERE plan = 'basic_yearly' OR plan = 'business_yearly'
