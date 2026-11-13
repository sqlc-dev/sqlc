-- name: PriceAsOf :many
SELECT id, amount FROM prices FOR SYSTEM_TIME AS OF @at WHERE id = @id;

-- name: PriceHistory :many
SELECT id, amount, valid_from, valid_to FROM prices FOR SYSTEM_TIME BETWEEN @from AND @to;
