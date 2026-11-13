-- name: Numeric :many
SELECT
  royalties * 2 AS doubled,
  royalties / units AS per_unit,
  units * weight AS weighted,
  ratio + units AS ratio_units,
  ratio * weight AS ratio_weight,
  units + total AS units_total,
  refunds - flags AS net_refunds,
  fee * units AS fee_units,
  small_fee + units AS small_fee_units,
  small_fee + fee AS fees,
  units % refunds AS remainder
FROM sales;

-- name: Concatenation :many
SELECT
  code + code AS codes,
  code + title AS code_title,
  ncode + title AS ncode_title,
  title + name AS title_name,
  code + ncode AS code_ncode
FROM sales;
