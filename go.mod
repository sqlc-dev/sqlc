module github.com/kyleconroy/sqlc

go 1.14

require (
	github.com/antlr/antlr4 v0.0.0-20200209180723-1177c0b58d07
	github.com/davecgh/go-spew v1.1.1
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/go-cmp v0.4.0
	github.com/jackc/pgx/v4 v4.6.0
	github.com/jinzhu/inflection v1.0.0
	github.com/lfittl/pg_query_go v1.0.0
	github.com/lib/pq v1.4.0
	github.com/pingcap/parser v0.0.0-20200623164729-3a18f1e5dceb
	github.com/spf13/cobra v1.0.0
	gopkg.in/yaml.v3 v3.0.0-20200121175148-a6ecf24a6d71
	vitess.io/vitess v0.0.0-20200617014457-5ba6549015c0
)

replace github.com/pingcap/parser => github.com/kyleconroy/parser v0.0.0-20200819185651-2caf0f596c0c
