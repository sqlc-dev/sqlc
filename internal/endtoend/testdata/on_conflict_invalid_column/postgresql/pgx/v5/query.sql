-- name: AddItem :exec
INSERT INTO cart_items (owner_id, product_id, price_amount, price_currency)
VALUES ($1, $2, $3, $4)
ON CONFLICT (owner_id, product_id) DO UPDATE
    SET price_amount1 = EXCLUDED.price_amount1,
        price_currency = EXCLUDED.price_currency;
