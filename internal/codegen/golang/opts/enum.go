package opts

import "fmt"

type SQLDriver string

const (
	SQLPackagePGXV4         string = "pgx/v4"
	SQLPackagePGXV5         string = "pgx/v5"
	SQLPackageStandard      string = "database/sql"
	SQLPackageYugaBytePGXV5 string = "yb/pgx/v5"
)

var validPackages = map[string]struct{}{
	string(SQLPackagePGXV4):         {},
	string(SQLPackagePGXV5):         {},
	string(SQLPackageYugaBytePGXV5): {},
	string(SQLPackageStandard):      {},
}

func validatePackage(sqlPackage string) error {
	if _, found := validPackages[sqlPackage]; !found {
		return fmt.Errorf("unknown SQL package: %s", sqlPackage)
	}
	return nil
}

const (
	SQLDriverPGXV4            SQLDriver = "github.com/jackc/pgx/v4"
	SQLDriverPGXV5                      = "github.com/jackc/pgx/v5"
	SQLDriverLibPQ                      = "github.com/lib/pq"
	SQLDriverGoSQLDriverMySQL           = "github.com/go-sql-driver/mysql"
	SQLDriverYugaBytePGXV5              = "github.com/yugabyte/pgx/v5"
)

var validDrivers = map[string]struct{}{
	string(SQLDriverPGXV4):            {},
	string(SQLDriverPGXV5):            {},
	string(SQLDriverLibPQ):            {},
	string(SQLDriverGoSQLDriverMySQL): {},
	string(SQLDriverYugaBytePGXV5):    {},
}

func validateDriver(sqlDriver string) error {
	if _, found := validDrivers[sqlDriver]; !found {
		return fmt.Errorf("unknown SQL driver: %s", sqlDriver)
	}
	return nil
}

func (d SQLDriver) IsPGX() bool {
	return d == SQLDriverPGXV4 || d == SQLDriverPGXV5 || d == SQLDriverYugaBytePGXV5
}

func (d SQLDriver) IsGoSQLDriverMySQL() bool {
	return d == SQLDriverGoSQLDriverMySQL
}

func (d SQLDriver) Package() string {
	switch d {
	case SQLDriverPGXV4:
		return SQLPackagePGXV4
	case SQLDriverPGXV5:
		return SQLPackagePGXV5
	case SQLDriverYugaBytePGXV5:
		return SQLPackageYugaBytePGXV5
	default:
		return SQLPackageStandard
	}
}
