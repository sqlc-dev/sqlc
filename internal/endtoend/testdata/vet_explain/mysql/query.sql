-- name: SelectById :one
SELECT id FROM debug
WHERE id = ? LIMIT 1;

-- name: SelectByCsmallint :one
SELECT id FROM debug
WHERE Csmallint = ? LIMIT 1;

-- name: SelectByCint :one
SELECT id FROM debug
WHERE Cint = ? LIMIT 1;

-- name: SelectByCinteger :one
SELECT id FROM debug
WHERE Cinteger = ? LIMIT 1;

-- name: SelectByCdecimal :one
SELECT id FROM debug
WHERE Cdecimal = ? LIMIT 1;

-- name: SelectByCnumeric :one
SELECT id FROM debug
WHERE Cnumeric = ? LIMIT 1;

-- name: SelectByCfloat :one
SELECT id FROM debug
WHERE Cfloat = ? LIMIT 1;

-- name: SelectByCreal :one
SELECT id FROM debug
WHERE Creal = ? LIMIT 1;

-- name: SelectByCdoubleprecision :one
SELECT id FROM debug
WHERE Cdoubleprecision = ? LIMIT 1;

-- name: SelectByCdouble :one
SELECT id FROM debug
WHERE Cdouble = ? LIMIT 1;

-- name: SelectByCdec :one
SELECT id FROM debug
WHERE Cdec = ? LIMIT 1;

-- name: SelectByCfixed :one
SELECT id FROM debug
WHERE Cfixed = ? LIMIT 1;

-- name: SelectByCtinyint :one
SELECT id FROM debug
WHERE Ctinyint = ? LIMIT 1;

-- name: SelectByCbool :one
SELECT id FROM debug
WHERE Cbool = ? LIMIT 1;

-- name: SelectByCmediumint :one
SELECT id FROM debug
WHERE Cmediumint = ? LIMIT 1;

-- name: SelectByCbit :one
SELECT id FROM debug
WHERE Cbit = ? LIMIT 1;

-- name: SelectByCdate :one
SELECT id FROM debug
WHERE Cdate = ? LIMIT 1;

-- name: SelectByCdatetime :one
SELECT id FROM debug
WHERE Cdatetime = ? LIMIT 1;

-- name: SelectByCtimestamp :one
SELECT id FROM debug
WHERE Ctimestamp = ? LIMIT 1;

-- name: SelectByCtime :one
SELECT id FROM debug
WHERE Ctime = ? LIMIT 1;

-- name: SelectByCyear :one
SELECT id FROM debug
WHERE Cyear = ? LIMIT 1;

-- name: SelectByCchar :one
SELECT id FROM debug
WHERE Cchar = ? LIMIT 1;

-- name: SelectByCvarchar :one
SELECT id FROM debug
WHERE Cvarchar = ? LIMIT 1;

-- name: SelectByCbinary :one
SELECT id FROM debug
WHERE Cbinary = ? LIMIT 1;

-- name: SelectByCvarbinary :one
SELECT id FROM debug
WHERE Cvarbinary = ? LIMIT 1;

-- name: SelectByCtinyblob :one
SELECT id FROM debug
WHERE Ctinyblob = ? LIMIT 1;

-- name: SelectByCblob :one
SELECT id FROM debug
WHERE Cblob = ? LIMIT 1;

-- name: SelectByCmediumblob :one
SELECT id FROM debug
WHERE Cmediumblob = ? LIMIT 1;

-- name: SelectByClongblob :one
SELECT id FROM debug
WHERE Clongblob = ? LIMIT 1;

-- name: SelectByCtinytext :one
SELECT id FROM debug
WHERE Ctinytext = ? LIMIT 1;

-- name: SelectByCtext :one
SELECT id FROM debug
WHERE Ctext = ? LIMIT 1;

-- name: SelectByCmediumtext :one
SELECT id FROM debug
WHERE Cmediumtext = ? LIMIT 1;

-- name: SelectByClongtext :one
SELECT id FROM debug
WHERE Clongtext = ? LIMIT 1;

-- name: SelectByCenum :one
SELECT id FROM debug
WHERE Cenum = ? LIMIT 1;

-- name: SelectByCset :one
SELECT id FROM debug
WHERE Cset = ? LIMIT 1;

-- name: SelectByCjson :one
SELECT id FROM debug
WHERE Cjson = ? LIMIT 1;

--
--

-- -- name: DeleteById :exec
-- DELETE FROM debug
-- WHERE id = ?;

-- -- name: DeleteByCsmallint :exec
-- DELETE FROM debug
-- WHERE Csmallint = ? LIMIT 1;

-- -- name: DeleteByCint :exec
-- DELETE FROM debug
-- WHERE Cint = ? LIMIT 1;

-- -- name: DeleteByCinteger :exec
-- DELETE FROM debug
-- WHERE Cinteger = ? LIMIT 1;

-- -- name: DeleteByCdecimal :exec
-- DELETE FROM debug
-- WHERE Cdecimal = ? LIMIT 1;

-- -- name: DeleteByCnumeric :exec
-- DELETE FROM debug
-- WHERE Cnumeric = ? LIMIT 1;

-- -- name: DeleteByCfloat :exec
-- DELETE FROM debug
-- WHERE Cfloat = ? LIMIT 1;

-- -- name: DeleteByCreal :exec
-- DELETE FROM debug
-- WHERE Creal = ? LIMIT 1;

-- -- name: DeleteByCdoubleprecision :exec
-- DELETE FROM debug
-- WHERE Cdoubleprecision = ? LIMIT 1;

-- -- name: DeleteByCdouble :exec
-- DELETE FROM debug
-- WHERE Cdouble = ? LIMIT 1;

-- -- name: DeleteByCdec :exec
-- DELETE FROM debug
-- WHERE Cdec = ? LIMIT 1;

-- -- name: DeleteByCfixed :exec
-- DELETE FROM debug
-- WHERE Cfixed = ? LIMIT 1;

-- -- name: DeleteByCtinyint :exec
-- DELETE FROM debug
-- WHERE Ctinyint = ? LIMIT 1;

-- -- name: DeleteByCbool :exec
-- DELETE FROM debug
-- WHERE Cbool = ? LIMIT 1;

-- -- name: DeleteByCmediumint :exec
-- DELETE FROM debug
-- WHERE Cmediumint = ? LIMIT 1;

-- -- name: DeleteByCbit :exec
-- DELETE FROM debug
-- WHERE Cbit = ? LIMIT 1;

-- -- name: DeleteByCdate :exec
-- DELETE FROM debug
-- WHERE Cdate = ? LIMIT 1;

-- -- name: DeleteByCdatetime :exec
-- DELETE FROM debug
-- WHERE Cdatetime = ? LIMIT 1;

-- -- name: DeleteByCtimestamp :exec
-- DELETE FROM debug
-- WHERE Ctimestamp = ? LIMIT 1;

-- -- name: DeleteByCtime :exec
-- DELETE FROM debug
-- WHERE Ctime = ? LIMIT 1;

-- -- name: DeleteByCyear :exec
-- DELETE FROM debug
-- WHERE Cyear = ? LIMIT 1;

-- -- name: DeleteByCchar :exec
-- DELETE FROM debug
-- WHERE Cchar = ? LIMIT 1;

-- -- name: DeleteByCvarchar :exec
-- DELETE FROM debug
-- WHERE Cvarchar = ?;

-- -- name: DeleteByCbinary :exec
-- DELETE FROM debug
-- WHERE Cbinary = ? LIMIT 1;

-- -- name: DeleteByCvarbinary :exec
-- DELETE FROM debug
-- WHERE Cvarbinary = ? LIMIT 1;

-- -- name: DeleteByCtinyblob :exec
-- DELETE FROM debug
-- WHERE Ctinyblob = ? LIMIT 1;

-- -- name: DeleteByCblob :exec
-- DELETE FROM debug
-- WHERE Cblob = ? LIMIT 1;

-- -- name: DeleteByCmediumblob :exec
-- DELETE FROM debug
-- WHERE Cmediumblob = ? LIMIT 1;

-- -- name: DeleteByClongblob :exec
-- DELETE FROM debug
-- WHERE Clongblob = ? LIMIT 1;

-- -- name: DeleteByCtinytext :exec
-- DELETE FROM debug
-- WHERE Ctinytext = ? LIMIT 1;

-- -- name: DeleteByCtext :exec
-- DELETE FROM debug
-- WHERE Ctext = ? LIMIT 1;

-- -- name: DeleteByCmediumtext :exec
-- DELETE FROM debug
-- WHERE Cmediumtext = ? LIMIT 1;

-- -- name: DeleteByClongtext :exec
-- DELETE FROM debug
-- WHERE Clongtext = ? LIMIT 1;

-- -- name: DeleteByCenum :exec
-- DELETE FROM debug
-- WHERE Cenum = ? LIMIT 1;

-- -- name: DeleteByCset :exec
-- DELETE FROM debug
-- WHERE Cset = ? LIMIT 1;

-- -- name: DeleteByCjson :exec
-- DELETE FROM debug
-- WHERE Cjson = ? LIMIT 1;
