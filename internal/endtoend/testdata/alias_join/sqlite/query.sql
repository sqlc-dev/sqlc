-- name: RetrieveAllWithDeprecation :many
select
  arn,
  name,
  aws_rds_databases.engine
from
  aws_rds_databases r
  natural join aws_rds_databases_engines
;

-- name: RetrieveAllWithDeprecationOther :many
select
  arn,
  name,
  r.engine
from
  aws_rds_databases r
  natural join aws_rds_databases_engines
;
