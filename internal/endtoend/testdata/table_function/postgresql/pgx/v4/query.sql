/* name: GetTransaction :many */
SELECT
	json_extract(transactions.data, '$.transaction.signatures[0]'),
	json_group_array(instructions.value)
FROM
  transactions, 
	json_each(json_extract(transactions.data, '$.transaction.message.instructions')) AS instructions
WHERE
	transactions.program_id = $1
	AND json_extract(transactions.data, '$.transaction.signatures[0]') > $2
	AND json_extract(json_extract(transactions.data, '$.transaction.message.accountKeys'), '$[' || json_extract(instructions.value, '$.programIdIndex') || ']') = transactions.program_id
GROUP BY transactions.id
LIMIT $3;