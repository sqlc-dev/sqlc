-- name: RollingStockGetOrCreateIDs :many
SELECT number, rolling_stock_id
FROM masterdata.rolling_stock_get_or_create_ids(
    inputs   => @inputs::masterdata.rolling_stock_number_input[],
    on_date  => @on_date::timestamptz,
    train_id => @train_id::bigint
);
