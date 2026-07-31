SELECT * FROM transactions
WHERE jsonb_extract_path(transactions.data, '$.transaction.signatures[0]') @> to_jsonb(sqlc.arg('data')::text);
