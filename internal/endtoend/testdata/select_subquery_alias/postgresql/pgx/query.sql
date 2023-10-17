-- name: FindWallets :many
select id, address, balance, total_balance from 
(
	select id, address, balance,
	  sum(balance) over (order by balance desc rows between unbounded preceding and current row) as total_balance,
	  sum(balance) over (order by balance desc rows between unbounded preceding and current row) - balance as last_balance
	from wallets
	where type=$1
) amounts
where amounts.last_balance < $2;
