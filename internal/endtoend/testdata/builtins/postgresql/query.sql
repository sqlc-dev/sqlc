-- name: Abs :one
SELECT abs(-17.4);
-- name: Cbrt :one
SELECT cbrt(27.0);
-- name: Ceil :one
SELECT ceil(-42.8);
-- name: Ceiling :one
SELECT ceiling(-95.3);
-- name: Degrees :one
SELECT degrees(0.5);
-- name: Div :one
SELECT div(9,4);
-- name: Exp :one
SELECT exp(1.0);
-- name: Floor :one
SELECT floor(-42.8);
-- name: Ln :one
SELECT ln(2.0);
-- name: Log :one
SELECT log(100.0);
-- name: Logs :one
SELECT log(2.0, 64.0);
-- name: Mod :one
SELECT mod(9,4);
-- name: Pi :one
SELECT pi();
-- name: Power :one
SELECT power(9.0, 3.0);
-- name: Radians :one
SELECT radians(45.0);
-- name: Round :one
SELECT round(42.4);
-- name: Rounds :one
SELECT round(42.4382, 2);
-- name: Scale :one
SELECT scale(8.41);
-- name: Sign :one
SELECT sign(-8.4);
-- name: Sqrt :one
SELECT sqrt(2.0);
-- name: Trunc :one
SELECT trunc(42.8);
-- name: Truncs :one
SELECT trunc(42.4382, 2);
-- name: WidthBucketNumerics :one
SELECT width_bucket(5.35, 0.024, 10.06, 5);
-- name: WidthBucketTimestamps :one
SELECT width_bucket(now(), array['yesterday', 'today', 'tomorrow']::timestamptz[]);
