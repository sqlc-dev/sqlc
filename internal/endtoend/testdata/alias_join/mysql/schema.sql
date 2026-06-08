CREATE TABLE IF NOT EXISTS aws_rds_databases (
  account_id VARCHAR(20) NOT NULL,
  region VARCHAR(20) NOT NULL,
  arn VARCHAR(20) NOT NULL,
  name TEXT NOT NULL,
  engine TEXT NOT NULL,
  engine_version TEXT NOT NULL,

  -- tags is a JSON object
  tags TEXT NOT NULL,

  UNIQUE (account_id, region, arn) 
);

CREATE TABLE IF NOT EXISTS aws_rds_databases_engines (
  engine VARCHAR(20) NOT NULL,
  engine_version VARCHAR(20) NOT NULL,
  deprecation TEXT NOT NULL,

  UNIQUE (engine, engine_version) 
);
