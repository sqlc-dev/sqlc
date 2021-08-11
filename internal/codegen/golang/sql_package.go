package golang

type SQLPackage string

const (
	SQLPackagePGX      SQLPackage = "pgx/v4"
	SQLPackageStandard SQLPackage = "database/sql"
)

func SQLPackageFromString(s string) SQLPackage {
	switch s {
	case string(SQLPackagePGX):
		return SQLPackagePGX
	default:
		return SQLPackageStandard
	}
}
