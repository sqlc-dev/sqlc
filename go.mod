module github.com/sqlc-dev/sqlc

go 1.26.0

require (
	github.com/antlr4-go/antlr/v4 v4.13.1
	github.com/cubicdaiya/gonp v1.0.4
	github.com/davecgh/go-spew v1.1.1
	github.com/fatih/structtag v1.2.0
	github.com/go-sql-driver/mysql v1.9.3
	github.com/google/cel-go v0.27.0
	github.com/google/go-cmp v0.7.0
	github.com/jackc/pgx/v4 v4.18.3
	github.com/jackc/pgx/v5 v5.8.0
	github.com/jinzhu/inflection v1.0.0
	github.com/lib/pq v1.11.2
	github.com/ncruces/go-sqlite3 v0.30.5
	github.com/pganalyze/pg_query_go/v6 v6.2.2
	github.com/pingcap/tidb/pkg/parser v0.0.0-20250324122243-d51e00e5bbf0
	github.com/riza-io/grpc-go v0.2.0
	github.com/spf13/cobra v1.10.2
	github.com/spf13/pflag v1.0.10
	github.com/sqlc-dev/doubleclick v1.0.0
	github.com/tetratelabs/wazero v1.11.0
	github.com/wasilibs/go-pgquery v0.0.0-20250409022910-10ac41983c07
	github.com/xeipuuv/gojsonschema v1.2.0
	golang.org/x/sync v0.19.0
	google.golang.org/grpc v1.79.1
	google.golang.org/protobuf v1.36.11
	gopkg.in/yaml.v3 v3.0.1
)

require (
	cel.dev/expr v0.25.1 // indirect
	filippo.io/edwards25519 v1.1.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.3 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgtype v1.14.0 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/ncruces/julianday v1.0.0 // indirect
	github.com/pingcap/errors v0.11.5-0.20240311024730-e056997136bb // indirect
	github.com/pingcap/failpoint v0.0.0-20240528011301-b51a646c7c86 // indirect
	github.com/pingcap/log v1.1.0 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	github.com/wasilibs/wazero-helpers v0.0.0-20240620070341-3dff1577cd52 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20180127040702-4e3ac2762d5f // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/crypto v0.47.0 // indirect
	golang.org/x/exp v0.0.0-20250620022241-b7579e27df2b // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)

replace github.com/go-sql-driver/mysql => github.com/sqlc-dev/mysql v0.0.0-20251129233104-d81e1cac6db2
