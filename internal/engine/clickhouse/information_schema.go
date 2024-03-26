package clickhouse

import "github.com/sqlc-dev/sqlc/internal/sql/catalog"

func getInformationSchema() *catalog.Schema {
	s := &catalog.Schema{Name: "information_schema"}
	s.Funcs = make([]*catalog.Function, 0)
	s.Tables = []*catalog.Table{
		//TODO
	}
	return s
}
